from datetime import datetime, timedelta
import io
import json
import requests
import string
import settings
import sys
from operator import itemgetter
from pprint import pprint
from performanceplatform.client import DataSet

'''
TODO:
- Write remaining tests
- Update when the PP adds pagination
- Currently we need to remove all data from the info-statistics dataset
  on the PP each time before running this script, via a curl command.
  So figure out how to do this conveniently, or how to store data
  usefully across time periods.
'''


class InfoStatistics():

    date_format = "%Y-%m-%dT00:00:00Z"
    date_format_longer = "%Y-%m-%dT00:00:00+00:00"

    def __init__(self):
        self.end_date = datetime.now()
        self.start_date = self.end_date - timedelta(days=settings.DAYS)

    def get_all_two_letter_prefixes(self):
        '''
        We need to iterate over the page-statistics dataset
        in the PP rather than getting it all at once, because
        the dataset is too large for the PP to return all at
        once, and the PP doesn't support pagination.
        So get all two-letter prefixes - we'll use these as
        the minimum viable way to iterate.
        '''
        prefixes = []
        letters = string.lowercase
        for letter1 in letters:
            for letter2 in letters:
                    prefixes.append('%s%s' % (letter1, letter2))
        return prefixes

    def construct_pp_query(self, dataset, value,
                           filter_by=None,
                           filter_by_prefix=None):
        '''
        Construct the PP URL to request.
        '''
        query = dataset
        query += '?group_by=pagePath&period=day'
        query += '&start_at=%s' % self.start_date.strftime(self.date_format)
        query += '&end_at=%s' % self.end_date.strftime(self.date_format)
        query += '&collect=%s' % value
        if filter_by_prefix:
            query += '&filter_by_prefix=pagePath:/%s' % filter_by_prefix
        if filter_by:
            query += '&filter_by=pagePath:%s' % filter_by
        return query

    def get_pp_data(self, query):
        '''
        Make PP call.
        '''
        results = []
        url = "%s/%s/%s" % (settings.DATA_DOMAIN, settings.DATA_GROUP, query)
        try:
            r = requests.get(url)
            if r.status_code == 200:
                json_data = r.json()
                if 'data' in json_data:
                    return json_data['data']
            else:
                print r.status_code, url
        except requests.exceptions.ConnectionError, requests.exceptions.HTTPError:
            print 'ERROR', url
        return []

    # def tidy_page_contacts(self, dataset):
    #     for d in dataset:

    #     return dataset

    def get_missing_page_statistics(self, output):
        '''
        Hopefully we should be able to delete this in
        six weeks' time when there's no longer lots of
        rubbish in the PP dataset.
        '''
        count = 0
        for o in output:
            if not o['uniquePageviews']:
                count += 1
                query = self.construct_pp_query('page-statistics',
                                                'uniquePageviews:sum',
                                                filter_by=o['pagePath'])
                print query
                data = self.get_pp_data(query)
                if len(data):
                    o['uniquePageviews'] = data[0]['uniquePageviews:sum']
        return output

    def get_urls_with_useful_metrics(self, results):
        '''
        Our dataset should only contain URLs that have some problem
        reports or searches. (This is because the page views
        dataset currently contains many meaningless, 404 URLs,
        which we want to exclude from our statistics.)
        '''
        useful_urls = [r['pagePath'] for r in results['page-contacts']]
        useful_urls += [r['pagePath'] for r in results['search-terms']]
        return sorted(set(useful_urls))

    def initialise_results(self, urls):
        '''
        For our set of unique URLs, initialise the data we will send
        to the PP: it needs certain datefields.
        '''
        output = []
        fields = ['uniquePageviews', 'problemReports',
                  'problemsPer100kViews', 'problemsNormalised',
                  'problemsQuintile', 'searchUniques',
                  'searchesPer100kViews', 'searchesNormalised',
                  'searchesQuintile', 'format', 'title']
        for url in urls:
            r = {}
            r['pagePath'] = url
            r['_timestamp'] = self.start_date.strftime(self.date_format)
            r['_start_at'] = self.start_date.strftime(self.date_format)
            r['_end_at'] = self.end_date.strftime(self.date_format)
            for f in fields:
                r[f] = None
            output.append(r)
        return output

    def get_title_and_format(self, output):
        '''
        Relies on a local file with all URLs returned by the Search API.
        I added this because I wanted to add formats to the entries
        in the one-off run for info-statistics. We need to look at a more
        sustainable way of doing this, or else
        '''

        json_data = open('./results/all-urls.json')
        data = json.load(json_data)
        all_urls = {}
        for d in data:
            link = d['link']
            all_urls[link] = {}
            all_urls[link]['format'] = d['format']
            all_urls[link]['title'] = d['title']

        for o in output:
            path = o['pagePath']
            if path in all_urls:
                o['format'] = all_urls[path]['format']
                o['title'] = all_urls[path]['title']
        return output

    def get_values_per_dataset(self, output, results):
        '''
        For our set of unique URLs, get the pageviews, problem
        reports and searches for each URL, from our PP datasets.
        '''
        for d in settings.DATASETS:
            dataset_name = d['dataset']
            value = d['value']
            outputValue = d['nice_value']
            ppValues = results[d['dataset']]
            for p in ppValues:
                for o in output:
                    if o['pagePath'] == p['pagePath']:
                        o[outputValue] = p[value]
                        if o[outputValue]:
                            o[outputValue] = float(o[outputValue])

        print 'Sorting smart answers...'
        # Sum problem reports for smart answers
        smart_answers = []
        for o in output:
            if o['format'] == 'smart-answer':
                smart_answers.append(o['pagePath'])
                # Add problem reports anything matching the start of the URL.
                problemReports = results['page-contacts']
                for p in problemReports:
                    if p['pagePath'].startswith(o['pagePath']):
                        if p['total:sum'] and o['problemReports']:
                            o['problemReports'] += float(p['total:sum'])
        return output

    def normalise_values(self, output):
        '''
        Calculate problem reports and searches per 100k views.
        Also calculate an experimental, normalised value.
        '''
        for r in output:
            if r['uniquePageviews']:
                if r['problemReports']:
                    raw = r['problemReports'] / r['uniquePageviews']
                    if r['problemReports'] > 5:
                        r['problemsPer100kViews'] = raw * 100000
                        r['problemsNormalised'] = r['problemsPer100kViews'] * r['problemReports']
                    else:
                        r['problemsPer100kViews'] = None
                        r['problemsNormalised'] = None
                if r['searchUniques']:
                    raw = r['searchUniques'] / r['uniquePageviews']
                    if r['searchUniques'] > 5:
                        r['searchesPer100kViews'] = raw * 100000
                        r['searchesNormalised'] = r['searchesPer100kViews'] * r['searchUniques']
                    else:
                        r['searchesPer100kViews'] = None
                        r['searchesNormalised'] = None
        return output

    def calculate_quintiles(self, output):
        '''
        TODO: Calculate quintile values for a given dataset
        Sort the dataset and then calculate quintile values.
        '''
        num_items = len(output)
        quintile_bound = num_items / 5.0
        output_by_problems = sorted(output,
                                    key=itemgetter('problemsPer100kViews'))
        for i, o in enumerate(output_by_problems):
            quintile = i // quintile_bound
            if o['problemsPer100kViews']:
                o['problemsQuintile'] = quintile + 1
        output_by_searches = sorted(output_by_problems,
                                    key=itemgetter('searchesPer100kViews'))
        for i, o in enumerate(output_by_searches):
            quintile = i // quintile_bound
            if o['searchesPer100kViews']:
                o['searchesQuintile'] = quintile + 1
        return output_by_searches

    def process_data(self):
        '''
        Iterate over letters of the alphabet, since the PP API doesn't have
        pagination, and will 503 if you ask for the full dataset.
        '''

        # Get PP values for all three datasets.
        # TODO: Error handling.
        print 'Getting data from PP...'
        results = {}
        prefixes = self.get_all_two_letter_prefixes()
        for d in settings.DATASETS:
            print d['dataset']
            results[d['dataset']] = []
            for prefix in prefixes:
                # if prefix != 'pi':
                #     continue
                print prefix
                query = self.construct_pp_query(d['dataset'],
                                                d['value'],
                                                filter_by_prefix=prefix)
                results[d['dataset']] += self.get_pp_data(query)

        # Get useful URLs.
        print 'Processing data...'
        urls = self.get_urls_with_useful_metrics(results)
        print len(urls), 'useful URLs found'
        output = self.initialise_results(urls)
        # Only run this if you have a local copy of all URLs
        # from the search API.
        # TODO: Remove.
        # output = self.get_title_and_format(output)

        # Fill in the missing page statistics.
        print 'Getting missing data...'
        output = self.get_values_per_dataset(output, results)
        print 'Getting missing page statistics'
        output = self.get_missing_page_statistics(output)

        # Do calculations.
        print 'Normalising values...'
        output = self.normalise_values(output)
        print 'Calculating quintiles...'
        output = self.calculate_quintiles(output)
        # print output[0]
        # pprint(output)

        # Dump results to a local JSON file.
        # This stores results in case posting fails: it also
        # gives us a historical archive of results.
        # TODO: Discuss whether this is useful/necessary.
        print 'Writing to local file...'
        fname = './results/data-%s.json' % self.end_date.strftime("%Y-%m-%d")
        with io.open(fname, 'w', encoding='utf-8') as f:
            f.write(unicode(json.dumps(output, indent=2, sort_keys=True)))

        # Empty PP dataset, then post new results. TODO: Error handling.
        print 'Posting data to PP...'
        data_set = DataSet.from_group_and_type(settings.DATA_DOMAIN,
                                               settings.DATA_GROUP,
                                               settings.RESULTS_DATASET,
                                               token=settings.PP_TOKEN)
        # data_set.post([])
        data_set.post(output)

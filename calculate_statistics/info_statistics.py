import io
import json
import os
import requests
import settings
import string
from datetime import datetime, timedelta
from performanceplatform.client import DataSet

'''
InfoStatistics class: this generates the aggregated data for
the PP's info-statistics dataset. This is used to identify
pages with high numbers of problem reports and searches.
It does the following:
- Fetches data from the PP for all pages with problem reports
  or searches
- Aggregates problem reports for smart answers to the level of
  the starting URL, to aid comparison
- Initialises a neat output dataset
- For all URLs with problem reports or searches, fetches data
  on the number of unique page views
- Normalise problem reports / searches by the number of unique
  page views
- Write output to a local JSON file and to the PP
'''


class InfoStatistics():

    date_format = "%Y-%m-%dT00:00:00Z"
    date_format_longer = "%Y-%m-%dT00:00:00+00:00"

    def __init__(self, pp_token):
        self.end_date = datetime.now()
        self.start_date = self.end_date - timedelta(days=settings.DAYS)
        self.pp_token = pp_token

    def construct_pp_query(self, dataset, value,
                           filter_by=None, filter_by_prefix=None):
        '''
        Construct the PP URL to request.
        '''
        query = dataset
        query += '?group_by=pagePath&period=day'
        query += '&start_at=%s' % self.start_date.strftime(self.date_format)
        query += '&end_at=%s' % self.end_date.strftime(self.date_format)
        query += '&collect=%s' % value
        if filter_by:
            query += '&filter_by=pagePath:%s' % filter_by
        elif filter_by_prefix:
            query += '&filter_by_prefix=pagePath:%s' % filter_by_prefix
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

    def get_smart_answer_urls(self):
        '''
        Get all smart answers, from the Search API.
        '''
        smart_answers = []
        url = 'https://www.gov.uk/api/search.json?filter_format=smart-answer'
        url += '&start=0&count=1000&fields=link'
        try:
            r = requests.get(url)
            if r.status_code == 200:
                json_data = r.json()
                if 'results' in json_data:
                    for l in json_data['results']:
                        smart_answers.append(l['link'])
        except requests.exceptions.ConnectionError, requests.exceptions.HTTPError:
            print 'ERROR', url

        return smart_answers

    def tidy_smart_answers(self, smart_answers, results):
        '''
        Sum page contacts values for smart answers, then
        remove non-root smart answers URLs.
        '''
        smart_answer_totals = {}
        for s in smart_answers:
            smart_answer_totals[s] = 0
        for r in results:
            for s in smart_answers:
                if 'total:sum' in r and r['pagePath'].startswith(s):
                    smart_answer_totals[s] += r['total:sum']

        for r in results:
            for s in smart_answer_totals:
                if r['pagePath'] in s and 'total:sum' in r:
                    r['total:sum'] = smart_answer_totals[s]

        for result in results[:]:
            path = result.get("pagePath")
            if "total:sum" in result and any(path.startswith(sa) and path != sa for sa in smart_answers):
                results.remove(result)

        return results

    def get_unique_urls(self, results):
        '''
        Return the unique URLs in the results from the
        first two datasets.
        '''
        urls = []
        for r in results:
            if r['pagePath'] not in urls:
                urls.append(r['pagePath'])
        return urls

    def initialise_results(self, urls):
        '''
        For our set of unique URLs, initialise the data we will send
        to the PP: it needs certain datefields.
        '''
        output = []
        fields = ['uniquePageviews', 'problemReports',
                  'problemsPer100kViews', 'problemsNormalised',
                  'searchUniques', 'searchesPer100kViews',
                  'searchesNormalised']
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

    def combine_dataset_results(self, output, results):
        '''
        Combine the results from the page-contacts and
        search-terms datasets into one.
        '''
        for r in results:
            for o in output:
                if o['pagePath'] == r['pagePath']:
                    if 'total:sum' in r:
                        o['problemReports'] = r['total:sum']
                    if 'searchUniques:sum' in r:
                        o['searchUniques'] = r['searchUniques:sum']
        return output

    def calculate_normalised_values(self, output):
        '''
        Calculate problem reports + searches per 100k views.
        Ignore pages with more problem reports than page views,
        as a primitive spam filter.
        Also calculate a normalised value.
        '''
        for r in output:
            if r['problemReports'] and r['uniquePageviews']:
                raw = r['problemReports'] / r['uniquePageviews']
                if r['problemReports'] > 2 and (r['problemReports'] < r['uniquePageviews']):
                    r['problemsPer100kViews'] = raw * 100000
                    r['problemsNormalised'] = r['problemsPer100kViews'] * r['problemReports']
            if r['searchUniques'] and r['uniquePageviews']:
                raw = r['searchUniques'] / r['uniquePageviews']
                if r['searchUniques'] > 3:
                    r['searchesPer100kViews'] = raw * 100000
                    r['searchesNormalised'] = r['searchesPer100kViews'] * r['searchUniques']
        return output

    def process_data(self):
        '''
        Main function.
        '''

        print 'Fetching PP data on problems and searches...'
        results = []
        for s in string.lowercase:
            query = self.construct_pp_query('page-contacts', 'total:sum',
                                            filter_by=None,
                                            filter_by_prefix='/%s' % s)
            results += self.get_pp_data(query)
            query = self.construct_pp_query('search-terms',
                                            'searchUniques:sum',
                                            filter_by=None,
                                            filter_by_prefix='/%s' % s)
            results += self.get_pp_data(query)

        # Tidy up results for smart answers.
        # It would be nice to amend how Feedex sends results to the PP
        # so that this step isn't necessary.
        print 'Tidying problem data for smart answers...'
        smart_answers = self.get_smart_answer_urls()
        results = self.tidy_smart_answers(smart_answers, results)

        print 'Initialising results...'
        urls = self.get_unique_urls(results)
        output = self.initialise_results(urls)
        output = self.combine_dataset_results(output, results)

        # For each URL, get page-statistics results from the PP.
        print 'Getting page view data per URL...'
        for o in output:
            query = self.construct_pp_query('page-statistics',
                                            'uniquePageviews:sum',
                                            filter_by=o['pagePath'])
            d = self.get_pp_data(query)
            if d and d[0]['uniquePageviews:sum']:
                o['uniquePageviews'] = int(d[0]['uniquePageviews:sum'])

        # Calculate normalised values.
        print 'Normalising values...'
        output = self.calculate_normalised_values(output)

        # Dump results to a local JSON file.
        # This stores results in case posting fails: it also
        # gives us a historical archive of results.
        # TODO: Discuss whether this is useful/necessary.
        print 'Writing to local file...'
        fname = './results/data-%s.json' % self.end_date.strftime("%Y-%m-%d")
        with io.open(fname, 'w', encoding='utf-8') as f:
            f.write(unicode(json.dumps(output, indent=2, sort_keys=True)))

        # Empty PP dataset, then post new results.
        # TODO: Error handling: not sure how the Python library does this?
        print 'Posting data to PP...'
        data_set = DataSet.from_group_and_type(settings.DATA_DOMAIN,
                                               settings.DATA_GROUP,
                                               settings.RESULTS_DATASET,
                                               token=self.pp_token)
        data_set.post(output)

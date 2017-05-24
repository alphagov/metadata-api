import unittest
from datetime import datetime
import settings
from info_statistics import InfoStatistics


class TestInfoStatistics(unittest.TestCase):

    def testConstructPPQuery(self):
        """construct_pp_query should return expected querystring"""
        dataset = 'page-contacts'
        value = 'total:sum'
        search_term = 'aa'
        query = info.construct_pp_query(dataset, value, search_term)
        d = 'page-contacts?group_by=pagePath&period=day'
        d += '&start_at=2014-12-16T00:00:00Z'
        d += '&end_at=2015-01-27T00:00:00Z'
        d += '&collect=total:sum&filter_by=pagePath:aa'
        self.assertEqual(query, d)
        query = info.construct_pp_query(dataset, value, filter_by=None,
                                        filter_by_prefix=search_term)
        d = 'page-contacts?group_by=pagePath&period=day'
        d += '&start_at=2014-12-16T00:00:00Z'
        d += '&end_at=2015-01-27T00:00:00Z'
        d += '&collect=total:sum&filter_by_prefix=pagePath:aa'
        self.assertEqual(query, d)

    def testTidySmartAnswers(self):
        """tidy_smart_answers should take an array of raw data
        results from the PP, and for all URLs that are smart answers,
        it should sum the number of problem reports and remove
         all entries except those for root URLs"""
        smart_answers = [u'/calculate-state-pension', u'/check-uk-visa']
        data = [
            {
                "pagePath": "/check-uk-visa",
                "total:sum": 2.0
            },
            {
                "pagePath": "/check-uk-visa/y",
                "total:sum": 3.0
            },
            {
                "pagePath": "/check-uk-visa/n",
                "total:sum": 4.0
            },
            {
                "pagePath": "/bank-holidays",
                "total:sum": 2.0
            },
            {
                "pagePath": "/check-uk-visa",
                "searchUniques:sum": 2.0
            }
        ]
        results = info.tidy_smart_answers(smart_answers, data)
        self.assertEqual(len(results), 3)
        self.assertEqual(results[0]['pagePath'], "/check-uk-visa")
        self.assertEqual(results[0]['total:sum'], 9.0)
        self.assertEqual(results[1]['pagePath'], "/bank-holidays")
        self.assertEqual(results[1]['total:sum'], 2.0)
        self.assertEqual(results[2]['pagePath'], "/check-uk-visa")
        self.assertEqual(results[2]['searchUniques:sum'], 2.0)

    def testGetUniqueURLs(self):
        """get_unique_urls should return a list of unique urls
        from an input dataset"""
        data = [
            {
                "pagePath": "/driving-licence-fees",
                "total:sum": 10.0
            },
            {
                "pagePath": "/driving-licence-fees/y",
                "total:sum": 10.0
            },
            {
                "pagePath": "/driving-licence-fees",
                "total:sum": 10.0
            }]
        results = info.get_unique_urls(data)
        expected = ["/driving-licence-fees", "/driving-licence-fees/y"]
        self.assertEqual(results, expected)

    def testInitialiseResults(self):
        """initialise_results should take a list of URLs and
        return an initialised set of results"""
        urls = ['/dartford-crossing',
                '/dartford-crossing-fees',
                '/dartford-crossing-fees-exemptions-penalties']
        results = info.initialise_results(urls)
        self.assertEqual(len(results), 3)
        self.assertEqual(results[0]['pagePath'], '/dartford-crossing')
        self.assertEqual(results[0]['_start_at'], '2014-12-16T00:00:00Z')
        self.assertEqual(results[0]['_timestamp'], '2014-12-16T00:00:00Z')
        self.assertEqual(results[0]['_end_at'], '2015-01-27T00:00:00Z')
        self.assertEqual(results[0]['searchUniques'], None)
        self.assertEqual(results[0]['searchesNormalised'], None)

    def testCombineDatasetResults(self):
        """combine_dataset_results should take an initialised
        set of results, and our raw data from the PP, and combine
        results from the raw data into a tidy set of output"""
        results = [
            {
                "pagePath": "/check-uk-visa",
                "total:sum": 9.0
            },
            {
                "pagePath": "/check-uk-visa",
                "searchUniques:sum": 2.0
            },
            {
                "pagePath": "/bank-holidays",
                "total:sum": 2.0
            }
        ]
        output = [
            {
                "pagePath": "/check-uk-visa",
                "problemReports": None,
                "searchUniques": None
            },
            {
                "pagePath": "/bank-holidays",
                "problemReports": None,
                "searchUniques": None
            }
        ]
        actual = info.combine_dataset_results(output, results)
        self.assertEqual(actual[0]['pagePath'], "/check-uk-visa")
        self.assertEqual(actual[0]['problemReports'], 9.0)
        self.assertEqual(actual[0]['searchUniques'], 2.0)
        self.assertEqual(actual[1]['pagePath'], "/bank-holidays")
        self.assertEqual(actual[1]['problemReports'], 2.0)
        self.assertEqual(actual[1]['searchUniques'], None)

    def testCalculateNormalisedValues(self):
        """calculate_normalise_values should return problems and searches
        per 100k page views, and normalised results, and deal with
        null values """
        data = [
            {
                'pagePath': u'/calculate-state-pension',
                'problemReports': 469.0,
                'searchUniques': 930.0,
                'uniquePageviews': 10000.0,
                'problemsPer100kViews': None,
                'problemsNormalised': None,
                'searchesPer100kViews': None,
                'searchesNormalised': None

            },
            {
                'pagePath': u'/calculate-employee-redundancy-pay',
                'problemReports': 20.0,
                'searchUniques': None,
                'uniquePageviews': None,
                'problemsPer100kViews': None,
                'problemsNormalised': None,
                'searchesPer100kViews': None,
                'searchesNormalised': None
            }
        ]
        results = info.calculate_normalised_values(data)
        self.assertEqual(results[0]['problemsPer100kViews'], 4690.0)
        self.assertEqual(results[0]['problemsNormalised'], 2199610.0)
        self.assertEqual(results[0]['searchesPer100kViews'], 9300.0)
        self.assertEqual(results[0]['searchesNormalised'], 8649000.0)
        self.assertEqual(results[1]['problemsPer100kViews'], None)
        self.assertEqual(results[1]['problemsNormalised'], None)
        self.assertEqual(results[1]['searchesPer100kViews'], None)
        self.assertEqual(results[1]['searchesNormalised'], None)


if __name__ == "__main__":
    info = InfoStatistics('foo')
    info.start_date = datetime.strptime("2014-12-16", "%Y-%m-%d").date()
    info.end_date = datetime.strptime("2015-01-27", "%Y-%m-%d").date()
    unittest.main()

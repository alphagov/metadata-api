import unittest
from datetime import datetime
import settings
from info_statistics import InfoStatistics

'''
Unit tests for InfoStatistics.
TODO: Finish stubbed tests, add missing tests
TODO: Do I need end-to-end tests?
'''


class KnownValues(unittest.TestCase):

    def testGetTwoLetterPrefixes(self):
        """get_all_two_letter_prefixes should do just that"""
        prefixes = info.get_all_two_letter_prefixes()
        self.assertEqual(prefixes[0], 'aa')
        self.assertEqual(prefixes[4], 'ae')
        self.assertEqual(len(prefixes), 26*26)

    def testConstructPPQuery(self):
        """construct_pp_query should return expected querystring"""
        dataset = 'page-contacts'
        value = 'total:sum'
        search_term = 'aa'
        query = info.construct_pp_query(dataset, self.start_date,
                                        self.end_date,
                                        value, search_term)
        d = 'page-contacts?group_by=pagePath&period=day'
        d += '&start_at=2014-12-16T00:00:00Z'
        d += '&end_at=2015-01-27T00:00:00Z'
        d += '&collect=total:sum&filter_by_prefix=pagePath:/aa'
        self.assertEqual(query, d)

    def testGetURLsWithUsefulMetrics(self):
        """get_urls_with_useful_metrics should return only URLs in the
        page-contacts and search-terms datasets"""
        results = {
            'page-statistics': [
                {'pagePath': '/dartford-crossing'},
                {'pagePath': '/dartford-crossing-fees'},
                {'pagePath': '/dartford-crossing-fees-exemptions-penalties'}
            ],
            'page-contacts': [
                {'pagePath': '/dartford-crossing'}
            ],
            'search-terms': [
                {'pagePath': '/dartford-crossing-fees'}
            ]
        }
        urls = info.get_urls_with_useful_metrics(results)
        self.assertEqual(len(urls), 2)
        self.assertEqual(urls[0], '/dartford-crossing')
        self.assertEqual(urls[1], '/dartford-crossing-fees')

    def testInitialiseResults(self):
        """initialise_results should take a list of URLs and
        return an initialised set of results"""
        urls = ['/dartford-crossing',
                '/dartford-crossing-fees',
                '/dartford-crossing-fees-exemptions-penalties']
        results = info.initialise_results(urls, self.start_date, self.end_date)
        self.assertEqual(len(results), 3)
        self.assertEqual(results[0]['pagePath'], '/dartford-crossing')
        self.assertEqual(results[0]['_start_at'], '2014-12-16T00:00:00Z')
        self.assertEqual(results[0]['_timestamp'], '2014-12-16T00:00:00Z')
        self.assertEqual(results[0]['_end_at'], '2015-01-27T00:00:00Z')
        self.assertEqual(results[0]['searchUniques'], None)
        self.assertEqual(results[0]['searchesNormalised'], None)

    def testGetValuesPerDataset(self):
        """get_values_per_dataset should..."""
        pass

    def testNormaliseValues(self):
        """normalise_values should..."""
        pass

    def testCalculateQuintiles(self):
        """calculate_quintiles should..."""
        pass

if __name__ == "__main__":
    info = InfoStatistics()
    info.start_date = datetime.strptime("2014-12-16", "%Y-%m-%d").date()
    info.end_date = datetime.strptime("2015-01-27", "%Y-%m-%d").date()
    unittest.main()

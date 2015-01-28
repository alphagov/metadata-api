DATA_DOMAIN = 'https://www.performance.service.gov.uk/data'
DATA_GROUP = 'govuk-info'
DATASETS = [
    {
        'dataset': 'page-statistics',
        'value': 'uniquePageviews:sum',
        'nice_value': 'uniquePageviews'
    },
    {
        'dataset': 'page-contacts',
        'value': 'total:sum',
        'nice_value': 'problemReports'
    },
    {
        'dataset': 'search-terms',
        'value': 'searchUniques:sum',
        'nice_value': 'searchUniques'
    }
]
DAYS = 42
RESULTS_DATASET = 'info-statistics'

# Add token here to run the script.
PP_TOKEN = ''

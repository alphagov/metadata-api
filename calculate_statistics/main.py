import os
import info_statistics
import settings

if 'PP_DATASET_TOKEN' not in os.environ:
    msg = 'You need to set the dataset token for the PP '
    msg += '%s/%s ' % (settings.DATA_GROUP, settings.RESULTS_DATASET)
    msg += 'dataset to run this script. You can get this from '
    msg += 'https://stagecraft.production.performance.service.gov.uk/admin/'
    print msg
else:
    c = info_statistics.InfoStatistics(os.environ['PP_DATASET_TOKEN'])
    c.process_data()

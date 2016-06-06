package performanceclient

/*
Package performanceclient is a client for the Performance Platform APIs.

Typical use is to read data from the Performance Platform.

To retrieve information about the dashboards:

    client := NewMetaClient("https://some.host/", logrus.New())
    dashboards, err := client.FetchDashboards()


To retrieve data for a particular module:
    
    client := NewDataClient("https://some.host/", logrus.New())
    res, err := client.Fetch("govuk-info", "page-statistics", QueryParams{})
    err = json.Unmarshall(res.Data, &v)
    
*/

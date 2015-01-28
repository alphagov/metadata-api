Work in progress!

Script to load aggregate statistics for info pages into the Performance
Platform's /info-statistics dataset.

This will be queried by metadata-api and returned to /info pages, in order
to put the numbers displayed on those pages into context.

It could also be used to power a dashboard of the pages with the most
searches, most problem reports etc.

Instructions
------------

Before loading the data, add the token for the PP dataset to your settings.
(NB: this should be replaced by an environment variable from `pp-puppet`.)

Install the dependencies (you may want to do this inside a virtualenv):

    pip install -r requirements.txt

NB: This should eventually be done in the project makefile.

Run the script to load data:

    python main.py

Testing
-------

Run tests:

    python test.py

NB: These should eventually be added to the project makefile.

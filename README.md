# An app

# TODO

### Main
* database back ups
    * currently cron-ed command but, should be servable /backup with auth header
* loop db connection attempt until success
    * this should actually happen on the DB.Ping call
    * will allow removal of init.sh
* better error handling than panic() on most errors
    * use log package

### Future plans
* find a way to run this service alongside regular website
    * create type which satisfies the HTTPHandler interface and stores the length of the prefix to the URL and acts as a semi-forward proxy to the in package http handlers

README

1. Anyone who wants use this repo will need Postgres and Go installed on their system to run the program.

2. Users will also have to install the gator CLI using:

"go install github.com/open-policy-agent/gatekeeper/v3/cmd/gator@master"

3. You must also create and save this file to your ROOT directory (not the project ROOT):

".gatorconfig.json"

which must contain "{"db_url":"postgres://{user name}:{password}@localhost:5432/gator?sslmode=disable"

4. Once everything is installed/updated correctly, you may run the following commands:
    1.  login
    2.  register
    3.  reset
    4.  users
    5.  agg (+time delay, if any, i.e. "agg 30s") 
    6.  addfeed
    7.  feeds
    8.  follow
    9.  following
    10. unfollow
    11. browse

using "go run . {command}" syntax


Thanks everyone for reading this and if you decide to try it out I hope you enjoy it!!

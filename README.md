# Betera-test-task

## API Routes
> /apods  

**GET** /apods?date=YYYY-MM-DD  
    Returns array of all APODS from album. Specify _date_ query param to get APODS for specific date.

**POST** /apods  
    Add APOD to album. Requires *File-Name* header and body should be .png, .jpeg or .gif image


## Steps to launch
1. Clone source code from github `git clone https://github.com/Hermes-Bird/betera-test-task`
2. Sign up at [Dropbox](https://www.dropbox.com/)
3. Go to <https://www.dropbox.com/developers> and create an app 
   with `sharing.write` & `files.content.write` permisions
4. Go to the app console and generate access token
5. Create `.env` file and paste access token along with other 
   config information (See `.env.example`)
6. Run `make run-docker-app`

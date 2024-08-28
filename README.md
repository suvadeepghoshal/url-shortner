# URL Shortner

### Topics intended:

- CRUD functionalities
- Database integration
- Error handling
- Logging

### Description:

erreate a URL shortener service where users can shorten long URLs and get a short URL that redirects to the original URL.
Track the number of times each short URL is accessed.

### Steps to Start:

- Set up a database to store original and shortened URLs.
- Create HTML templates for inputting and displaying URLs.
- Implement the logic for generating and resolving short URLs.
- Track and log URL access events.

### Need to Do:
- [X] Session Management
- [X] User Auth
- [X] Short URL Logic completion
- [X] Long URL extraction and redirection
- [ ] Do DB calls in separate Go routines
- [ ] Add logging service
- [ ] Add common response handler
- [ ] Add common error handler
- [ ] Linking short URLs to users
- [ ] Implement cache to retrieve URLs for a user

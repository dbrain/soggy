# Work In Progress
  This is a work in progress of an (express)[https://github.com/visionmedia/express/] like web .. thing .. for Go. Mostly a "Danny wants to play with Go". Use it at your own risk, if it ever becomes something pretty I'll let you know, Brian (I really hope your name is Brian).

# What it does currently
  * Allows you to add routes
  * Mount servers to a path
  * Allow you to add middleware
  * Probably breaks, first time playing with Go, most likely using pointers incorrectly.

# Stuff
The commander-express.go has some kind of example of how to use this. It's what I'm running to see I haven't broken anything.

# Testing
Ok, some tests were added so it's not so HAHA worthy. But they're really basic and probably lacking coverage.

# Developer rambling to himself goes here
TODO:
- Look at the loop in soggy.go, can probably be changed now that servers handle the requests and nexting. Basically finding first true and running a command on it.

Mission List:
* Parameters in URL
* Get the template engines going, server.Engine('.html', MagicHtmlEngine) where MagicHtmlEngine implements a specific interface. Response.Render('file', params) to render (file can be relative to a specific path, or use server.Config vars)
* Environment specific config (move the app.use etc. into a configure function, maybe)
* Error Handling
* Body Parser
* Steal cool things from everyone else, then spread nasty rumours about their mothers.
* Look into this whole FCGI / SCGI thing.

Error Handling:
- Get this working instead of just panicing

Validation:
- Look at how other frameworks do validation, work out how that would fit into express.go

Body parser:
- Parse JSON into map[string]interface{} or a specified struct.. how to specify?

Head requests:
- Maybe head requests should go into any routes and if its head it doesn't perform the writes, just the head .. or something ?

How to do Next:
- Fake loop using nexts
- Call the first middleware passing next function (that calls the next middleware using index)

TODO:
- Expand the normal http.Request and ResponseWriter into their own structs.. so helper params / functions can be added
- Work out if im using pointers vs values correctly ? Reponse/Request mostly
- Need to put next function back in, otherwise it's impossible to have middleware that survives from the start to the end of a request. Work out how to do this well


- End of the day all request handlers (middleware, routes etc.) should be added to the one array
- Request comes in, loops middleware executing all of them
- Might need to make next smart to pass stuff previous middleware has collected, won't just be able to add junk to request
- Maybe make next take err and a map of string interface{}.. shove stuff on there.
- Route middleware should have a single handler function which picks its route

Handle:
// Is mounting worthwhile at all? Guess the handler could chain servers?
- Allow mounting, so one server can bind /api to a module that acts like a seperate app (so it can be pulled out later)
- Make createServer return something you can pass to http.Server as a handle -- app.RequestHandler is a HandlerFunc.. maybe add app.Handler() to return a handler should be simple enough (ignore until important)
- express.createServer().listen(3000) should function (listening on 0.0.0.0:3000)

- Allow for binding render handlers to filetypes via express like API e.g. exports.__express = function(filename, options, callback) through express.createServer().engine('.html', magicTemplateEngine)
// - express.createServer().use(func (req, res, next) || func (err, req, res, next)) middleware
// - express.createServer().use(express.createServer.Router) router middleware
// - Add middleware through funcs get/head/post/put/delete
- Allow and parse parameters out of the url
- express.createServer().set('config option', 'value')
- Environment specific config
- Body parser middleware that will write to a parsedBody (JSON unmarshalled etc.)
- File uploads need to work or at least make room for them


func Handler(req, res, next) {

}

func ErrorHandler(err, req, res, next) {

}

func PlayingWithHeaders(...) {
  res.Set('Content-Type', 'application/json')
  res.Get('Content-Type')
}

func PlayingWithReq(...) {
  req.Accepts('application/json')
  req.Protocol // read from X-Forwarded-Proto then server
  req.AcceptsCharset('charset')
  req.AcceptsLanguage('lang')
  req.Stale / req.Fresh // to see if a request is stale (based on ETag/Last-Modified)
}

func PlayingWithRes(...) {
  // Code is always first (or 200 if not specified (doable in go?))
  res.Send(code, body) // write and close
  res.Write('BLAH') // can write multiple times
  res.End() // finish writing
  res.Redirect(code, url)
  res.Json(code, body) // write and close with Content-Type as application/json, maybe marshal as well depending on body type
  res.Render('template_ref', optsForTemplate(map[string]interface{}), callback)
  res.Type('json') // set Content-Type to json with good charset
  res.cookies // handle this however, JSON cookie support would be cool
}

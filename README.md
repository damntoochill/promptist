# Promptist

Welcome to the repo for Promptist. A place for AI artists to share, chat, and learn prompt craft. We want to create an online "third place" for the AI art scene.

## Open Source Platform, r u mad?

We wanted to try things a little bit differently with building the platform. Instead of placing the project under lock and key, we are going to let the community have a say in
the development of  the platform. The goal is to have all features be determined by the community.

## Tech Stack

### Go

We are using Go on the backend with MySQL. We are also using Go for server side rendering. So there is no javascript frontend at the moment.

### Tailwind CSS
We are using tailwind for CSS design.

```
tailwindcss -i ./assets/css/input.css -o ./public/css/output.css --watch
```

## Discord

Come join us in [Discord](https://discord.com/invite/WcXPatW9sY).


## Dev

Download the Go code, set some env vars, and run it.

### Env vars
```
PORT=
APP_NAME=
DB_DATABASE=
DB_USER=
DB_PASSWORD=
DB_HOST=
DB_PORT=
SESSION_KEY=
MAILGUN_DOMAIN=
MAILGUN_API_KEY=
MAILGUN_PUBLIC_API_KEY=
EMAIL_SENDER=
WEBSITE_HOST=
S3_KEY=
S3_SECRET=
S3_ENDPOINT=
S3_REGION=
S3_BUCKET=
```

### Useful tools

- [air](https://github.com/cosmtrek/air) - live reload for the web app
- [direnv](https://direnv.net) - load env variables easily

### Ugly Code

I didn't intend to open source this project, so the code might be a little ugly until we can get to a more
standardized place.
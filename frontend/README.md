## Build the Docker Image:
`docker build -t my-nextjs-app .`

## Run the Docker Container
`docker run -p 3000:3000 my-nextjs-app`

## For testing
go into dashboard/page.js there is comment add to comment that line
after that run `npx jest --verbose`
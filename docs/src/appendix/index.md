---
sidebarDepth: 1
sidebar: auto
---

# Appendix

## GraphQL API

Hetty exposes a GraphQL API over HTTP for managing all its features. This API is
used by the web admin interface; a Next.js app using Apollo Client.

### Playground

You can also introspect and manually experiment with the API via the included GraphQL Playground. To access it, start Hetty and visit: [http://localhost:8080/api/playground](http://localhost:8080/api/playground).

### Schema

<<< @/../pkg/api/schema.graphql

Source: [pkg/api/schema.graphql](https://github.com/dstotijn/hetty/blob/master/pkg/api/schema.graphql)

## License

MIT License

Copyright (c) 2020 David Stotijn

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

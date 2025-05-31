# Canvus API Endpoint Table

This table summarizes all Canvus API endpoints, with resource, method, path, key parameters, and a short description. Every resource and endpoint is fully enumerated—no omissions or placeholders.

| Resource         | Method | Path                                              | Key Parameters                        | Description                                      |
|------------------|--------|--------------------------------------------------|----------------------------------------|--------------------------------------------------|
| Canvases         | GET    | /canvases                                        | subscribe (query)                      | List all canvases                               |
| Canvases         | GET    | /canvases/:id                                    | id (path), subscribe (query)           | Get a single canvas                            |
| Canvases         | POST   | /canvases                                        | name, folder_id (body)                 | Create a canvas                                |
| Canvases         | PATCH  | /canvases/:id                                    | id (path), name, mode (body)           | Update (rename/mode) a canvas                  |
| Canvases         | DELETE | /canvases/:id                                    | id (path)                              | Delete a canvas                                |
| Canvases         | GET    | /canvases/:id/preview                            | id (path)                              | Get canvas preview                             |
| Canvases         | POST   | /canvases/:id/restore                            | id (path)                              | Restore demo canvas                            |
| Canvases         | POST   | /canvases/:id/save                               | id (path)                              | Save demo state                                |
| Canvases         | POST   | /canvases/:id/move                               | id, folder_id (path/body)              | Move/trash a canvas                            |
| Canvases         | POST   | /canvases/:id/copy                               | id, folder_id (path/body)              | Copy a canvas                                  |
| Canvases         | GET    | /canvases/:id/permissions                        | id (path)                              | Get permissions                                |
| Canvases         | POST   | /canvases/:id/permissions                        | id (path), permissions (body)          | Set permissions                                |
| Canvas Folders   | GET    | /canvas-folders                                  | subscribe (query)                      | List all folders                               |
| Canvas Folders   | GET    | /canvas-folders/:id                              | id (path), subscribe (query)           | Get a single folder                            |
| Canvas Folders   | POST   | /canvas-folders                                  | name, folder_id (body)                 | Create a folder                                |
| Canvas Folders   | PATCH  | /canvas-folders/:id                              | id (path), name (body)                 | Rename a folder                                |
| Canvas Folders   | POST   | /canvas-folders/:id/move                         | id, folder_id (path/body)              | Move a folder to another folder or to trash     |
| Canvas Folders   | POST   | /canvas-folders/:id/copy                         | id, folder_id (path/body)              | Copy a folder inside another folder             |
| Canvas Folders   | DELETE | /canvas-folders/:id                              | id (path)                              | Delete a folder                                |
| Canvas Folders   | DELETE | /canvas-folders/:id/children                     | id (path)                              | Delete all children of a folder                 |
| Canvas Folders   | GET    | /canvas-folders/:id/permissions                  | id (path)                              | Get permission overrides on a folder            |
| Canvas Folders   | POST   | /canvas-folders/:id/permissions                  | id (path), permissions (body)          | Set permission overrides on a folder            |
| Notes            | GET    | /canvases/:id/notes                              | id (path)                              | List all notes of the specified canvas          |
| Notes            | GET    | /canvases/:id/notes/:note_id                     | id, note_id (path)                     | Get a single note                              |
| Notes            | POST   | /canvases/:id/notes                              | id (path), note (body)                 | Create a note                                  |
| Notes            | PATCH  | /canvases/:id/notes/:note_id                     | id, note_id (path), note (body)        | Update a note                                  |
| Notes            | DELETE | /canvases/:id/notes/:note_id                     | id, note_id (path)                     | Delete a note                                  |
| Images           | GET    | /canvases/:id/images                             | id (path)                              | List all images of the specified canvas         |
| Images           | GET    | /canvases/:id/images/:image_id                   | id, image_id (path)                    | Get a single image                             |
| Images           | POST   | /canvases/:id/images                             | id (path), image (multipart)           | Create an image                                |
| Images           | PATCH  | /canvases/:id/images/:image_id                   | id, image_id (path), image (body)      | Update an image                                |
| Images           | DELETE | /canvases/:id/images/:image_id                   | id, image_id (path)                    | Delete an image                                |
| PDFs             | GET    | /canvases/:id/pdfs                               | id (path)                              | List all PDFs of the specified canvas           |
| PDFs             | GET    | /canvases/:id/pdfs/:pdf_id                       | id, pdf_id (path)                      | Get a single PDF                               |
| PDFs             | GET    | /canvases/:id/pdfs/:pdf_id/download              | id, pdf_id (path)                      | Download a single PDF                          |
| PDFs             | POST   | /canvases/:id/pdfs                               | id (path), pdf (multipart)             | Create a PDF                                   |
| PDFs             | PATCH  | /canvases/:id/pdfs/:pdf_id                       | id, pdf_id (path), pdf (body)          | Update a PDF                                   |
| PDFs             | DELETE | /canvases/:id/pdfs/:pdf_id                       | id, pdf_id (path)                      | Delete a PDF                                   |
| Videos           | GET    | /canvases/:id/videos                             | id (path)                              | List all videos of the specified canvas         |
| Videos           | GET    | /canvases/:id/videos/:video_id                   | id, video_id (path)                    | Get a single video                             |
| Videos           | GET    | /canvases/:id/videos/:video_id/download          | id, video_id (path)                    | Download a single video                        |
| Videos           | POST   | /canvases/:id/videos                             | id (path), video (multipart)           | Create a video                                 |
| Videos           | PATCH  | /canvases/:id/videos/:video_id                   | id, video_id (path), video (body)      | Update a video                                 |
| Videos           | DELETE | /canvases/:id/videos/:video_id                   | id, video_id (path)                    | Delete a video                                 |
| Widgets          | GET    | /canvases/:id/widgets                            | id (path)                              | List all widgets of the specified canvas        |
| Widgets          | GET    | /canvases/:id/widgets/:widget_id                 | id, widget_id (path)                   | Get a single widget                            |
| Widgets          | POST   | /canvases/:id/widgets                            | id (path), widget (body)               | Create a widget                                |
| Widgets          | PATCH  | /canvases/:id/widgets/:widget_id                 | id, widget_id (path), widget (body)    | Update a widget                                |
| Widgets          | DELETE | /canvases/:id/widgets/:widget_id                 | id, widget_id (path)                   | Delete a widget                                |
| Anchors          | GET    | /canvases/:id/anchors                            | id (path)                              | List all anchors of the specified canvas        |
| Anchors          | GET    | /canvases/:id/anchors/:anchor_id                 | id, anchor_id (path)                   | Get a single anchor                            |
| Browsers         | GET    | /canvases/:id/browsers                           | id (path)                              | List all browsers of the specified canvas       |
| Browsers         | GET    | /canvases/:id/browsers/:browser_id               | id, browser_id (path)                  | Get a single browser                           |
| Connectors       | GET    | /canvases/:id/connectors                         | id (path)                              | List all connectors for the specified canvas    |
| Connectors       | GET    | /canvases/:id/connectors/:connector_id           | id, connector_id (path)                | Get a single connector                         |
| Connectors       | POST   | /canvases/:id/connectors                         | id (path), connector (body)            | Create a connector                             |
| Connectors       | PATCH  | /canvases/:id/connectors/:connector_id           | id, connector_id (path), connector (body)| Update a connector                          |
| Connectors       | DELETE | /canvases/:id/connectors/:connector_id           | id, connector_id (path)                | Delete a connector                             |
| Annotations      | GET    | /canvases/:canvasId/widgets?annotations=1        | canvasId (path), annotations (query)    | List all annotations for widgets on a canvas    |
| Annotations      | GET    | /canvases/:canvasId/widgets?annotations=1&subscribe=1 | canvasId (path), annotations, subscribe (query) | Subscribe to annotation changes (streaming) |
| Canvas Backgrounds| GET   | /canvases/:id/background                         | id (path)                              | Get canvas background                          |
| Canvas Backgrounds| PATCH | /canvases/:id/background                         | id (path), background (body)           | Set canvas background (solid color or haze)     |
| Canvas Backgrounds| POST  | /canvases/:id/background                         | id (path), image (multipart)           | Set canvas background to an image               |
| Color Presets    | GET    | /canvases/:canvasId/color-presets                | canvasId (path)                        | Get color presets for a canvas                  |
| Color Presets    | PATCH  | /canvases/:canvasId/color-presets                | canvasId (path), presets (body)        | Update color presets for a canvas               |
| Clients          | GET    | /clients                                        |                                          | List all clients connected to the server        |
| Clients          | GET    | /clients/:id                                    | id (path)                              | Get a single client                            |
| Groups           | GET    | /groups                                         |                                          | List all user groups                           |
| Groups           | GET    | /groups/:id                                     | id (path)                              | Get a single user group                        |
| Groups           | POST   | /groups                                         | name (body), description (body)         | Create a new user group                        |
| Groups           | DELETE | /groups/:id                                     | id (path)                              | Delete a user group                            |
| Groups           | POST   | /groups/:group_id/members                       | group_id (path), id (body)              | Add a user to a group                          |
| Groups           | GET    | /groups/:id/members                             | id (path)                              | List users in a group                          |
| Groups           | DELETE | /groups/:group_id/members/:user_id              | group_id, user_id (path)                | Remove a user from a group                     |
| Audit Log        | GET    | /audit-log                                      | filters (query)                         | Get the list of audit events (paginated)        |
| Audit Log        | GET    | /audit-log/export-csv                           | filters (query)                         | Export the audit log as CSV                     |
| License          | GET    | /license                                        |                                          | Get license info                               |
| License          | POST   | /license/activate                               | key (body)                              | Activate license online                         |
| License          | GET    | /license/request                                | key (query)                             | Generate offline activation request             |
| License          | POST   | /license                                        | license (body)                           | Install license from offline activation         |
| Mipmaps & Assets | GET    | /mipmaps/{publicHashHex}                        | publicHashHex (path), canvas-id (header)| Get mipmap info for an asset                    |
| Mipmaps & Assets | GET    | /mipmaps/{publicHashHex}/{level}                | publicHashHex, level (path), canvas-id (header)| Get a specific mipmap level image (WebP)  |
| Mipmaps & Assets | GET    | /assets/{publicHashHex}                         | publicHashHex (path), canvas-id (header)| Get asset file by hash                          |
| Video Inputs     | GET    | /canvases/:id/video-inputs                      | id (path)                               | List all video input widgets on a canvas        |
| Video Inputs     | POST   | /canvases/:id/video-inputs                      | id (path), widget (body)                | Create a video input widget                     |
| Video Inputs     | DELETE | /canvases/:id/video-inputs/:input_id            | id, input_id (path)                     | Delete a video input widget                     |
| Video Inputs     | GET    | /clients/:client_id/video-inputs                | client_id (path)                        | List video input sources on a client            |
| Video Outputs    | GET    | /clients/:client_id/video-outputs               | client_id (path)                        | List all video outputs for a client device      |
| Video Outputs    | PATCH  | /clients/:client_id/video-outputs/:index        | client_id, index (path), source/suspended (body)| Set video output source or suspend      |
| Video Outputs    | PATCH  | /canvases/:id/video-outputs/:output_id          | id, output_id (path), name/resolution (body)| Update a video output (name, resolution)   |
| Uploads Folder   | POST   | /canvases/:id/uploads-folder                    | id (path), json/data (multipart)        | Upload a note (multipart POST)                  |
| Uploads Folder   | POST   | /canvases/:id/uploads-folder                    | id (path), json/data (multipart)        | Upload a file asset (multipart POST)            |
| Workspaces       | GET    | /clients/:client_id/workspaces                  | client_id (path)                        | List all workspaces of a client                 |
| Workspaces       | GET    | /clients/:client_id/workspaces/:workspace_index | client_id, workspace_index (path)       | Get a single workspace                          |
| Workspaces       | PATCH  | /clients/:client_id/workspaces/:workspace_index | client_id, workspace_index (path), params (body)| Update workspace parameters              |
| Access Tokens    | GET    | /users/:id/access-tokens                        | id (path)                              | List access tokens for a user                   |
| Access Tokens    | GET    | /users/:id/access-tokens/:token-id              | id, token-id (path)                    | Get info about a single access token            |
| Access Tokens    | POST   | /users/:id/access-tokens                        | id (path), description (body)           | Create a new access token                       |
| Access Tokens    | PATCH  | /users/:id/access-tokens/:token-id              | id, token-id (path), description (body) | Change access token description                 |
| Access Tokens    | DELETE | /users/:id/access-tokens/:token-id              | id, token-id (path)                     | Delete (revoke) an access token                 |
| Server Config    | GET    | /server-config                                  |                                          | Get server settings                             |
| Server Config    | PATCH  | /server-config                                  | settings (body)                         | Change server settings                          |
| Server Config    | POST   | /server-config/send-test-email                  |                                          | Send a test email                               |
| Server Info      | GET    | /server-info                                    |                                          | Get server info                                 |

---

For a structured list of all endpoints, see [Canvus_API_Endpoint_List.md](Canvus_API_Endpoint_List.md) 
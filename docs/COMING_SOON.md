# Canvus REST API - Proposed Features (Not Yet Implemented)

This document lists proposed REST API features that have not yet been implemented in Canvus.

---

## Client REST API

### Accessing Canvas Widgets Through Client Workspace API

**Proposed Endpoint Pattern:**
```
GET /clients/:client-id/workspaces/:workspace-id/browsers/:browser-widget-uuid/
```

**Benefits:**
- Use coordinates in workspace coordinate system instead of canvas coordinate system
- Simplifies code for positioning widgets (e.g., centering on screen)
- Natural place to expose client-specific widget properties:
  - Video streaming (receiving) state
  - Video output state
- Natural place to expose client-specific widgets
- Natural way to connect anchors and workspace in presentation mode
- Easy way to put a widget on the screen without querying which canvas is opened

---

### IP Video Widgets API

**Proposed Endpoints:**

List IP video widgets:
```
GET /clients/:client-id/workspaces/:workspace-id/ip-videos
```

Create widget:
```
POST /clients/:client-id/workspaces/:workspace-id/ip-videos

{
    "url": "http://video.example.com/webrtc",
    "location": <location>
}
```

Change play/pause state:
```
PATCH /clients/:client-id/workspaces/:workspace-id/ip-videos/:video-id

{
    "playback_state": "play"
}
```

---

### RDP Widgets API

**Proposed Endpoints:**

List configured RDP servers:
```
GET /clients/:client-id/rdp-servers
```

Access RDP connections (widgets):
```
GET /clients/:client-id/workspaces/:workspace-id/rdp-connections
```

---

### Presentation Mode API

**Proposed Endpoints:**

List presentations:
```
GET /canvases/:id/presentations
```

Get single presentation:
```
GET /canvases/:id/presentations/:id
```

Get slides:
```
GET /canvases/:id/presentations/:id/slides
```

Add slide:
```
POST /canvases/:id/presentations/:id/slides

{
    "prev": "<slide-id>",
    "next": null,
    "content-id": "<widget-id>",
    "transition": "<transition-type>",
    "touch": true,
    "temporary-annotations": false,
    "focus": false,
    "autoplay": true
}
```

Remove slide:
```
DELETE /canvases/:id/presentations/:id/slides/:slide-id
```

Change slide properties:
```
PATCH /canvases/:id/presentations/:id/slides/:slide-id

{
    "content-id": "<widget-id>",
    "transition": "<transition-type>",
    "autoplay": false,
    "temporary-annotations": false
}
```

Reorder slides:
```
PATCH /canvases/:id/presentations/:id/slides/:slide-id

{
    "prev": "<slide-id-1>",
    "next": "<slide-id-2>"
}
```

**Client Presentation API:**

Get presentation status:
```
GET /clients/:client-id/workspaces/:workspace-id/presentation

{
    "presentation_state": "none" | "presentation" | "widget",
    "presentation_id": "<presentation-id>",
    "slide_id": "<slide-id>",
    "widget_id": "<widget-id>"
}
```

Start presentation:
```
POST /clients/:client-id/workspaces/:workspace-id/presentation/start

{
    "presentation_id": "<presentation-id>",
    "slide_id": "<slide-id>",
    "widget_id": "<widget-id>"
}
```

Stop presentation:
```
POST /clients/:client-id/workspaces/:workspace-id/presentation/stop
```

Navigate slides:
```
POST /clients/:client-id/workspaces/:workspace-id/presentation/next-slide
POST /clients/:client-id/workspaces/:workspace-id/presentation/prev-slide
POST /clients/:client-id/workspaces/:workspace-id/presentation/set-slide

{
    "slide_id": "<slide-id>"
}
```

---

### Video Streaming State Control

**Proposed Endpoints:**

Put browser/ip-stream/rdp widget into streaming state:
```
PATCH /clients/:client-id/workspaces/:workspace-id/browsers/:browser-id

{
    "streaming_state": "streamer"
}
```

Put receiver to play/pause state:
```
PATCH /clients/:client-id/workspaces/:workspace-id/browsers/:browser-id

{
    "playback_state": "play"
}
```

---

## Canvas REST API

### Tables API

API to manipulate tables on canvas.

**Status:** On hold pending future of "Tables" feature.

---

### Annotations API

**Requirements:**
- Draw annotations
- Import/export to SVG (and raster?) format
- MVP using efficient internal data format

**Considerations:**
- Annotations are not widgets
- Adding as widget children makes objects huge
- Separate endpoint would create synchronization issues

---

### Bookmarks API

**Status:** Not currently requested, but missing from API coverage.

---

### /children Endpoint (Read Operations)

**Current State:** `/canvas-folders/:id/children` only supports DELETE.

**Proposed:** Add GET support as syntax sugar to simplify folder operations.

---

### Single Endpoint for Folders and Canvases

**Problem:** Canvases and folders are accessed through different endpoints (`/canvases` and `/canvas-folders`), requiring API users to subscribe to two different streams to build folder tree representation.

**Proposed:** Unified endpoint for both.

---

### Export/Import Canvases

**Proposed Endpoints:**
```
POST /canvases/:id/export
POST /canvases/import
```

---

### Export/Import Folders

**Proposed Endpoints:**
```
POST /canvas-folders/:id/export
POST /canvas-folders/import
```

---

### Async Operation Endpoint

**Problem:** Long-running operations (export/import/copy/move folders) need progress tracking.

**Proposed:** Operations optionally provide `/async_status` mode to:
- Check status
- Monitor progress
- Cancel operations

---

### Canvas Participants API

**Purpose:** Canvas-specific resource to access information for:
- Avatar list visualization
- User list
- Viewports of other users (bottom-right corner UI in desktop client)

---

## Assets API Improvements

### Current Endpoint
```
GET /api/v1/assets/:hash
```

### Issues with Current Implementation

1. Cannot be cached based on URL (same URL might return different data when image changes)
2. Cannot be cached using headers (server doesn't support HEAD requests)
3. Two instances of same widget don't share data (downloaded twice)
4. Potential sync issue if image changes while sending download request
5. HTTP `<video>` tag videos are unseekable (no HTTP range request support)

### Proposed Improvements

1. Add appropriate HTTP headers for caching
2. Add support for HEAD requests
3. Add support for HTTP range requests
4. Change hash in responses to include full public asset hash instead of truncated one

---

## WebSocket API

### Current Issues

1. Node.js backend creates unnecessary buffering, extra processing, memory consumption, latency
2. Multiple HTTP connections per web client consume server resources
3. `/widgets` endpoint format not ideal for huge canvases (one huge JSON array)
4. Annotations don't fit well in `/widgets` model
5. Moving widgets requires many PATCH requests per second for real-time feel
6. No mechanism for web client to advertise viewport/client info
7. No way to cache `/widgets` (can be 10MB+)

### Proposed Solution

**Single WebSocket Connection:**
- Open directly from browser to server (skip node.js backend)
- Authenticate once, process messages with less overhead
- Use "ndjson" format optimized for minimal overhead
- Client subscribes to different things over single connection
- Subscribe only to single canvas for simplicity
- Client sends attribute changes and workspace/client info through same channel
- Cache id (canvas modified-timestamp) for efficient reconnection after tab suspension

**Suggested Libraries:**
- Go: Gorilla WebSockets
- JS: robust-websocket
- Avoid socket.io (unnecessary overhead)

**Example Communication:**

```
Client: {"msg": "subscribe-canvas", "canvas": "<uuid>"}
Server: {"msg": "canvas-begin", "canvas": "<uuid>"}
Server: {"id":"cf2ade8e-...", "location":{"x":0,"y":0}, ...}
Server: <other canvas items...>
Server: {"msg": "canvas-end", "canvas": "<uuid>", "cache-id": <cache-id>}
```

On change:
```
Server: {"id": "<note uuid>", "location": {"x": 123, "y": 234}, "cache-id": <new cache-id>}
```

On reconnect with valid cache:
```
Client: {"msg": "subscribe-canvas", "canvas": "<uuid>", "cache-id": <cache-id>}
Server: {"msg": "canvas-end", "canvas": "<uuid>", "cache-id": <cache-id>}
```

---

## Known Issues

### Colors Missing Alpha Channel

**Problem:** Colors are formatted as `#RRGGBB` without alpha channel, causing canvases with transparent notes to appear broken in web client.

**Solution:** Support `#RRGGBBAA` format.

---

## Reference Documents

These proposals reference additional detailed documents:
- Video Output REST API
- Video Inputs REST API
- Presentations REST API
- Annotations REST API
- Mipmap REST API
- Connectors REST API
- Canvas Background REST API
- Participants REST API

# Annotations API

## Overview
The Annotations API in MT Canvus is accessed via the widgets endpoint. Annotations are always listed as part of the widget properties, not as a separate resource. This design allows efficient batch processing and ensures the web client receives all annotations for a widget at once.

## Listing All Annotations on a Canvas

To list all annotations for widgets on a canvas, use:

```
GET /canvases/:canvasId/widgets?annotations=1
```

- `:canvasId` is the UUID of the canvas.
- The `annotations=1` query parameter includes all annotations for each widget in the response.

### Response Structure
The response is an array of widget objects. Each widget includes an `annotations` field, which is an array of annotation objects. If a widget has no annotations, the `annotations` field is present and empty.

Example:
```json
[
  {
    // ... regular widget properties ...
    "annotations": [
      {
        "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f", // Annotation UUID
        "page": 2,                // (optional) Page number, omitted if < 1 or undefined
        "depth": 0,
        "line_color": "#FF0000",
        "points": "..."           // base64 encoded Float32Array of spline control points
      }
      // ... more annotations ...
    ]
  },
  // ... more widgets ...
]
```

- If a widget has no annotations, it will have: `"annotations": []`
- No `parent_id` or `state` fields in the annotation objects in this response.

## Subscribing to All Annotations on a Canvas

To subscribe to annotation changes:

```
GET /canvases/:canvasId/widgets?annotations=1&subscribe=1
```

- The initial response is the same as above.
- When annotations change, updates are sent as individual annotation objects ("strokes"). Only changed fields are included in the update payload.
- Changed widgets do not include the `annotations` field after the initial response.

### Annotation Event Payloads

#### Annotation Added
```json
[
  {
    "widget_type": "Annotation",
    "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f",
    "state": "normal",
    // ... page, depth, line_color, points ...
  }
]
```

#### Annotation Deleted
```json
[
  {
    "widget_type": "Annotation",
    "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f",
    "state": "deleted"
  }
]
```

#### Annotation Path Changed
```json
[
  {
    "widget_type": "Annotation",
    "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f",
    "points": "..." // Full base64 encoded Float32Array
  }
]
```

#### Individual Bezier Node Changes
- The `points` array consists of 3D cubic Bezier nodes (see Luminous::BezierNode). The number of floats is always a multiple of 9. There are always at least two nodes (18 floats).
- Events map to `Valuable::AttributeEvent` types:

##### Node Inserted (`ELEMENT_INSERTED`)
```json
[
  {
    "widget_type": "Annotation",
    "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f",
    "points:insert": [3, "base64-encoded-9-floats"]
  }
]
```

##### Node Changed (`ELEMENT_CHANGED`)
```json
[
  {
    "widget_type": "Annotation",
    "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f",
    "points:change": [3, "base64-encoded-9-floats"]
  }
]
```

##### Node Erased (`ELEMENT_ERASED`)
```json
[
  {
    "widget_type": "Annotation",
    "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f",
    "points:erase": 3
  }
]
```

## Notes
- The web client does not need to group or collect annotations; all are provided in batch per widget.
- If a widget has no annotations, the empty `annotations` field ensures old annotations are cleared on restore.
- All IDs are UUIDs generated by the MCS.
- This documentation is the source of truth for the API. The implementation must match this specification.

## Endpoints

### List All Annotations on a Canvas
- **Endpoint:** `/v1/canvases/{canvasId}/widgets?annotations=1`
- **Method:** `GET`
- **Description:** Lists all widgets on a canvas, with each widget including an `annotations` array containing all its annotations.
- **Headers:**
  - `Private-Token: <API_KEY>`

#### Example Request
```http
GET /v1/canvases/123/widgets?annotations=1
Private-Token: <API_KEY>
```

#### Example Response
```json
[
  {
    "id": "widgetA",
    // ... other widget properties ...
    "annotations": [
      {
        "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f",
        "page": 2,
        "depth": 0,
        "line_color": "#FF0000",
        "points": "..." // base64 encoded Float32Array
      }
      // ... more annotations ...
    ]
  },
  // ... more widgets ...
]
```

### Subscribe to Annotation Changes
- **Endpoint:** `/v1/canvases/{canvasId}/widgets?annotations=1&subscribe=1`
- **Method:** `GET` (streaming)
- **Description:** Subscribes to real-time annotation changes for all widgets on a canvas. The initial response is the same as above. Subsequent updates are sent as events for changed annotations only.
- **Headers:**
  - `Private-Token: <API_KEY>`

#### Example Event Payloads
- Annotation Added:
```json
[
  {
    "widget_type": "Annotation",
    "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f",
    "state": "normal",
    // ... other fields ...
  }
]
```
- Annotation Deleted:
```json
[
  {
    "widget_type": "Annotation",
    "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f",
    "state": "deleted"
  }
]
```
- Annotation Path Changed:
```json
[
  {
    "widget_type": "Annotation",
    "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f",
    "points": "..." // base64 encoded Float32Array
  }
]
```
- Individual Bezier Node Changes (insert/change/erase):
```json
[
  {
    "widget_type": "Annotation",
    "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f",
    "points:insert": [3, "base64-encoded-9-floats"]
  }
]
[
  {
    "widget_type": "Annotation",
    "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f",
    "points:change": [3, "base64-encoded-9-floats"]
  }
]
[
  {
    "widget_type": "Annotation",
    "id": "b7e8c2e2-1f2a-4c3d-9e5f-2a6b7c8d9e0f",
    "points:erase": 3
  }
]
```

### Creating, Updating, and Deleting Annotations
- **Annotations are not managed via a dedicated /annotations endpoint.**
- To create, update, or delete annotations, use the widgets API and the appropriate client actions (e.g., via the UI or SDK). The backend will emit the correct annotation events as shown above.
- There is no RESTful POST/PATCH/DELETE for `/annotations` as a separate resource.

---

**Note:**
- All annotation operations are performed as part of widget operations.
- The `annotations` field is always present (empty if none) to ensure clients can clear old annotations.
- All IDs are UUIDs generated by the server.
- This documentation reflects the actual, supported API. If you need to create or modify annotations, use the widget mechanisms provided by the client SDK or UI, not a direct REST endpoint.

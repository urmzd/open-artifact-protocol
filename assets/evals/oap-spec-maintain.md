## OAP Target Markers

For non-JSON text artifacts, major blocks and individually-updatable values are wrapped with target markers:

```
<oap:target id="ID">content</oap:target>
```

Targets nest — coarse blocks contain fine-grained value targets:

```html
<oap:target id="stats">
  <div class="card">
    <h3>Revenue</h3>
    <span><oap:target id="revenue-value">$12,340</oap:target></span>
  </div>
</oap:target>
```

Target IDs describe the role, not the current value (e.g., "total-revenue" not "12345").
Place targets where values are most likely to be revised.

For `application/json`, artifacts do not contain OAP marker tags. JSON edits use JSON Pointer paths such as `/data/66/email`.

## OAP Edit Envelope

To edit an artifact, produce a JSON envelope with `name: "edit"`:

```json
{
  "protocol": "oap/0.1",
  "id": "artifact-id",
  "version": 2,
  "name": "edit",
  "meta": {"format": "text/html"},
  "content": [
    {"op": "replace", "target": {"type": "id", "value": "revenue-value"}, "content": "$15,720"}
  ]
}
```

For non-JSON artifacts, target by ID only: `{"type": "id", "value": "target-id"}`. Reference existing target IDs from the artifact.

For JSON artifacts, target by pointer only: `{"type": "pointer", "value": "/path/to/value"}`. Replacement and insertion content must be a serialized JSON value.

Ops: `replace`, `delete`, `insert_before`, `insert_after`.

IMPORTANT: You MUST respond with a JSON edit envelope, NOT the full artifact. Reference only existing IDs from `list_targets()` for non-JSON artifacts or existing paths from `list_paths()` for JSON artifacts. Do not invent target IDs or paths. Use `replace` to update content within a target, `delete` to remove content at a target/path, and `insert_before`/`insert_after` to insert inside the selected target range or around the selected JSON array item. Always increment the version number.

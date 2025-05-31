# TODO

- If remote touch is enabled on an input window it stops being shown as a video-input in the API; this seems like a bug.
- We need to add floating / not floating option to the API for video-input windows.
- We need to add mute/unmute for video-inputs and videos endpoints.
- We need to change host_id in video-inputs to be client_id since this is the same data element, and having two names for the same thing is a challenge.
- Add API support for creating VideoOutputAnchor widgets.
- Fix bug to allow changing the video output source to a widget via API - currently gives a message that this widget is not supported but manually switching in the UI allows it to work.
- Complete or implement the annotations endpoint: Current API does not support creating or updating annotations via PATCH/POST, and GET on notes does not return annotations even with ?annotations=1. 
Investigate and implement full annotation support in the API.
# Protocol: Question -> Answer
__Based on UDP__

__Packet Structure:__
```
[Question ID (8 bytes)] [Repeats Count (1 byte)] [QA Marker (1 byte)]
```

__Question ID:__
- Randomic bytes, server returns answer with this ID.

__Repeats Count:__
- Count of resends datagram for this session and this Question ID. Set by client.

__QA Marker:__
- 0 - 64: This packet is question
- 191 - 255: This packet is answer

__What essence?__
Send question to server until a answer is received or the number of attempts is exhausted.

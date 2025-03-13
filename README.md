# schoolsubscriber
## time ranges file
A file with time ranges should be provided with the **-c** flag:
```yaml
ranges:
  - start: 2025-03-13 09:00:00
    end: 2025-03-14 00:20:00
  - start: 2025-03-14 15:00:00
    end: 2025-03-14 16:00:00
```

- The date-time format should be exactly like this.
- Multiple time ranges might be provided.
- The closest currently available slot from any of the ranges is occupied.

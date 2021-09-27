# 42-Vox


## Position Shifts
**Vertex Position :**
```go
y << 23 // mask 0x1FF (511)
x << 15 // mask 0xFF  (255)
z << 7  // mask 0xFF  (255)
```
## FaceData Shifts
**Textures Coordinates :**

```go 
x << 26 // mask 0x3F (63)
y << 20 // mask 0x3F (63)
```

**Billow Color Modifier :**
```go
c << 12 // mask 0xFF (255)
```

**Texture ID :**
```go
id << 5 // mask 0x7F (127)
```

**Face Type :**
```go
face << 2 // mask 0x7 (7)
```


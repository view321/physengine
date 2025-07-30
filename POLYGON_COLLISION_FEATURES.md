# Polygon Collision System

This document describes the polygon collision system added to the physics engine.

## Overview

The polygon collision system allows for collision detection and response between arbitrary polygon shapes, circles, and AABB boxes. It uses the Separating Axis Theorem (SAT) for accurate collision detection between polygons.

## Features

### 1. Polygon Collider Component

The `PolygonCollider` component stores a list of vertices that define the polygon shape:

```go
type PolygonColliderData struct {
    Vertices []Vec2.Vec2 // Local vertices relative to center
}
```

### 2. Collision Detection Types

The system supports the following collision combinations:

- **Polygon vs Polygon**: Uses SAT algorithm for accurate collision detection
- **Polygon vs Circle**: Finds the closest point on the polygon to the circle center
- **Polygon vs AABB**: Converts AABB to polygon vertices and uses SAT
- **Circle vs Circle**: Existing circle collision (unaffected)
- **AABB vs AABB**: Existing AABB collision (unaffected)

### 3. Helper Functions

#### Creating Regular Polygons

```go
// Create a regular polygon with n sides and given radius
vertices := components.CreateRegularPolygon(6, 50) // Hexagon with radius 50

// Create a rectangle with given width and height
vertices := components.CreateRectangle(120, 80) // Rectangle 120x80
```

#### Getting World Vertices

```go
// Get vertices transformed to world space (with rotation and position)
worldVertices := components.GetWorldVertices(entity)
```

### 4. Factory Functions

The system includes factory functions for creating different polygon shapes:

- `CreateTestTriangle()` - Creates a triangular polygon
- `CreateTestPentagon()` - Creates a pentagonal polygon  
- `CreateTestHexagon()` - Creates a hexagonal polygon
- `CreateTestRectangle()` - Creates a rectangular polygon
- `CreatePolygonCollisionDemo()` - Creates a complete demo scene

## Usage Examples

### Creating a Triangle

```go
entity := ecs.World.Create(components.MaterialComponent, components.Transform, 
    components.PolygonCollider, components.Drawable, components.MassComponent, 
    components.Velocity, components.AngularVelocity, components.Torque)
entry := ecs.World.Entry(entity)

// Set up transform and physics components
components.SetPos(entry, Vec2.Vec2{X: 200, Y: 150})
components.Velocity.Get(entry).Velocity = Vec2.Vec2{X: 80, Y: 20}

// Create triangle vertices
vertices := []Vec2.Vec2{
    {X: 0, Y: -50},   // top
    {X: -40, Y: 30},  // bottom-left
    {X: 40, Y: 30},   // bottom-right
}

poly := components.PolygonCollider.Get(entry)
poly.Vertices = vertices
```

### Creating a Regular Polygon

```go
// Create a hexagon with radius 50
vertices := components.CreateRegularPolygon(6, 50)
poly := components.PolygonCollider.Get(entry)
poly.Vertices = vertices
```

### Creating a Rectangle

```go
// Create a rectangle 120x80
vertices := components.CreateRectangle(120, 80)
poly := components.PolygonCollider.Get(entry)
poly.Vertices = vertices
```

## Collision Detection Algorithm

### Separating Axis Theorem (SAT)

The polygon collision detection uses the SAT algorithm:

1. **Project both polygons onto each axis** (normals of edges)
2. **Check for separation** - if projections don't overlap on any axis, no collision
3. **Find minimum overlap** - the axis with minimum overlap determines collision normal
4. **Calculate penetration depth** for collision response

### Circle vs Polygon

For circle-polygon collision:

1. **Find closest point** on polygon to circle center
2. **Calculate distance** from circle center to closest point
3. **Check if distance < radius** for collision
4. **Calculate collision normal** from circle center to closest point

## Physics Integration

The polygon collision system integrates with the existing physics system:

- **Angular impulse** calculation for rotational response
- **Friction** handling with static and dynamic friction
- **Positional correction** to prevent objects from sticking together
- **Material properties** (restitution, friction) affect collision response

## Performance Considerations

- **SAT algorithm** is O(n*m) where n and m are the number of vertices
- **Regular polygons** are more efficient than irregular ones
- **Convex polygons** are required for accurate collision detection
- **Vertex count** should be kept reasonable (typically < 20 vertices)

## Demo Scene

The `CreatePolygonCollisionDemo()` function creates a comprehensive demo scene with:

- Different polygon shapes (triangle, pentagon, hexagon, rectangle)
- Circles interacting with polygons
- AABB boxes interacting with polygons
- Stationary objects for collision testing
- Various velocities and angular velocities

## Future Enhancements

Potential improvements to the polygon collision system:

1. **Convex hull generation** for automatic polygon optimization
2. **Concave polygon support** with decomposition
3. **Broad phase collision detection** for performance optimization
4. **Polygon simplification** for complex shapes
5. **Custom collision shapes** with user-defined vertices 
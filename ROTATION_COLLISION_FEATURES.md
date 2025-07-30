# Rotation-Aware Collision System

This document describes the enhanced collision system that supports rotation-aware collision detection and response.

## Overview

The rotation-aware collision system extends the existing physics engine to properly handle collisions between rotated objects, including:

- **Rotated AABB vs Rotated AABB** collision detection using Separating Axis Theorem (SAT)
- **Circle vs Rotated AABB** collision detection
- **Angular impulse** in collision response
- **Rotational friction** effects

## New Components

### Rotated Collision Detection (`components/rotated_collision.go`)

#### `RotatedAABBvsAABB(a1, a2 *donburi.Entry) (bool, Vec2.Vec2, float64)`
- **Purpose**: Detects collision between two rotated AABBs
- **Returns**: (colliding, normal, penetration)
- **Algorithm**: Separating Axis Theorem (SAT)

#### `RotatedCirclevsAABB(circle, box *donburi.Entry) (bool, Vec2.Vec2, float64)`
- **Purpose**: Detects collision between a circle and a rotated AABB
- **Returns**: (colliding, normal, penetration)
- **Algorithm**: Transforms circle to box's local space, finds closest point

## New Systems

### `UpdateRotatedCollisions` (`systems/rotated_collisions.go`)

Replaces the original collision system with rotation-aware collision detection and response.

#### Collision Types Handled:
1. **Circle vs Circle** - Standard circle collision (rotation doesn't affect circles)
2. **Rotated AABB vs Rotated AABB** - SAT-based collision detection
3. **Circle vs Rotated AABB** - Transformed coordinate space collision

### `ResolveWithAngularImpulse` Function

Handles collision response with both linear and angular impulse:

```go
func ResolveWithAngularImpulse(e1, e2 *donburi.Entry, normal Vec2.Vec2, collisionPoint Vec2.Vec2, res1, res2 float64) float64
```

#### Physics Equations:
- **Linear Impulse**: `J = -(1 + e) * v_rel / (1/m1 + 1/m2 + (r1×n)²/I1 + (r2×n)²/I2)`
- **Angular Impulse**: `Δω = J * (r × n) / I`
- **Tangential Velocity**: `v_tangential = ω × r`

### `ResolveRotatedFriction` Function

Handles friction with rotational effects:

```go
func ResolveRotatedFriction(e1, e2 *donburi.Entry, normal Vec2.Vec2, collisionPoint Vec2.Vec2, j float64)
```

## Collision Detection Algorithms

### Separating Axis Theorem (SAT)

Used for rotated AABB collision detection:

1. **Get Axes**: Calculate normal vectors of all edges
2. **Project Polygons**: Project both polygons onto each axis
3. **Check Separation**: If projections don't overlap on any axis, no collision
4. **Find Minimum Overlap**: Determine collision normal and penetration depth

### Circle vs Rotated AABB

1. **Transform to Local Space**: Rotate circle center to box's coordinate system
2. **Find Closest Point**: Clamp circle center to box bounds
3. **Calculate Distance**: Distance from circle center to closest point
4. **Transform Normal**: Rotate collision normal back to world space

## Physics Integration

### Angular Impulse Calculation

The collision response now includes angular effects:

```go
// Calculate tangential velocities at collision point
tangentialVel1 := Vec2.Vec2{
    X: -angVel1.AngularVelocity * r1.Y,
    Y: angVel1.AngularVelocity * r1.X,
}

// Total velocity at collision point
velAtPoint1 := vel1.Velocity.Add(tangentialVel1)
```

### Impulse Denominator

The impulse calculation includes angular terms:

```go
denominator := m1.InverseMass + m2.InverseMass

// Add angular terms
if angVel1 != nil && m1.InverseInertia > 0 {
    cross1 := r1.X*normal.Y - r1.Y*normal.X
    denominator += cross1 * cross1 * m1.InverseInertia
}
```

## Demo Objects

### `CreateRotatingCollisionDemo`

Creates a set of objects to demonstrate rotation-aware collisions:

1. **Rotating Square**: Moves left with 3.0 rad/s rotation
2. **Rotating Circle**: Moves right with -2.0 rad/s rotation  
3. **Stationary Object**: Heavy object for collision testing

### Factory Functions

- `CreateRotatingSquare(ecs, pos, vel, angularVel)` - Creates rotating square
- `CreateRotatingCircle(ecs, pos, vel, angularVel)` - Creates rotating circle
- `CreateStationaryObject(ecs, pos, vel)` - Creates stationary heavy object

## System Order

The enhanced collision system runs in this order:

1. `UpdateCamera` - Updates camera position
2. `UpdateRotatedCollisions` - **NEW**: Rotation-aware collision detection and response
3. `UpdateVelocity` - Updates linear velocity
4. `ApplyContinuousTorque` - Applies continuous torque (demo)
5. `UpdateTorque` - Converts torque to angular velocity changes
6. `UpdateAngularVelocity` - Updates rotation based on angular velocity

## Performance Considerations

### SAT Algorithm Complexity
- **Time**: O(n*m) where n, m are the number of edges in each polygon
- **Space**: O(n+m) for storing axes and projections

### Optimization Features
- **Early Exit**: SAT stops when separation is found on any axis
- **Efficient Projections**: Only projects onto edge normals
- **Cached Calculations**: Reuses transformed coordinates

## Usage Examples

### Creating Objects with Rotation-Aware Collision

```go
// Create a rotating square with collision
entity := ecs.World.Create(
    components.Transform, 
    components.AABB_Component, 
    components.Velocity, 
    components.AngularVelocity, 
    components.MassComponent,
    components.MaterialComponent
)

// Set rotation and collision properties
components.SetAngularVelocity(entry, 2.0) // 2 rad/s rotation
components.SetPos(entry, Vec2.Vec2{X: 100, Y: 100})

// Configure collision bounds
aabb := components.AABB_Component.Get(entry)
aabb.Min = Vec2.Vec2{X: -50, Y: -50}
aabb.Max = Vec2.Vec2{X: 50, Y: 50}
```

### Collision Response with Angular Effects

The system automatically handles:
- **Linear impulse** based on collision normal
- **Angular impulse** based on collision point and moment arm
- **Friction** with both linear and angular components
- **Positional correction** to prevent overlap

## Benefits

1. **Realistic Physics**: Objects rotate realistically during collisions
2. **Accurate Collision Detection**: Proper handling of rotated shapes
3. **Consistent Response**: Angular and linear physics work together
4. **Performance**: Efficient algorithms for real-time simulation
5. **Extensible**: Easy to add new collision shapes and responses

This creates a complete rotation-aware collision system that provides realistic physics simulation for rotating objects. 
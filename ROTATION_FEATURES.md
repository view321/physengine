# Rotation Features in Physics Engine

This document describes the rotation features that have been added to the physics engine.

## New Components

### AngularVelocity Component
- **File**: `components/angular_velocity.go`
- **Purpose**: Stores the angular velocity (rotation speed) of an entity
- **Fields**: 
  - `AngularVelocity float64` - The angular velocity in radians per second

### Torque Component
- **File**: `components/torque.go`
- **Purpose**: Stores the torque (rotational force) applied to an entity
- **Fields**:
  - `Torque float64` - The torque in arbitrary units

## New Systems

### UpdateAngularVelocity System
- **File**: `systems/angular_velocity.go`
- **Purpose**: Updates entity rotation based on angular velocity and delta time
- **Formula**: `rotationDelta = angularVelocity * deltaTime`

### UpdateTorque System
- **File**: `systems/torque.go`
- **Purpose**: Applies torque to update angular velocity based on physics
- **Formula**: `angularAcceleration = torque / inertia`
- **Formula**: `angularVelocityDelta = angularAcceleration * deltaTime`

### ApplyContinuousTorque System
- **File**: `systems/continuous_torque.go`
- **Purpose**: Demonstrates continuous torque application for testing

## Enhanced Components

### MassComponent
- **Enhanced Fields**:
  - `Inertia float64` - Moment of inertia
  - `InverseInertia float64` - Inverse moment of inertia (1/inertia)

### Transform Component
- **Enhanced Functions**:
  - `GetRotationMatrix(angle float64)` - Returns 2x2 rotation matrix
  - `RotatePoint(point Vec2.Vec2, angle float64)` - Rotates a point around origin
  - `RotatePointAround(point Vec2.Vec2, center Vec2.Vec2, angle float64)` - Rotates a point around a center

## Helper Functions

### AngularVelocity Component
- `SetAngularVelocity(entry *donburi.Entry, angularVelocity float64)` - Sets angular velocity
- `ChangeAngularVelocity(entry *donburi.Entry, angularVelocityDelta float64)` - Changes angular velocity
- `GetAngularVelocity(entry *donburi.Entry) float64` - Gets current angular velocity

### Torque Component
- `SetTorque(entry *donburi.Entry, torque float64)` - Sets torque
- `AddTorque(entry *donburi.Entry, torque float64)` - Adds to current torque
- `GetTorque(entry *donburi.Entry) float64` - Gets current torque

## Usage Examples

### Creating a Rotating Object
```go
// Create entity with rotation components
entity := ecs.World.Create(components.Transform, components.AngularVelocity, components.Torque, components.MassComponent)
entry := ecs.World.Entry(entity)

// Set initial angular velocity
components.SetAngularVelocity(entry, 2.0) // 2 radians per second

// Apply torque
components.SetTorque(entry, 1000.0) // Apply torque

// Set mass and inertia
mc := components.MassComponent.Get(entry)
mc.Mass = 10
mc.Inertia = mc.Mass * radius * radius * 0.5 // For circle
mc.InverseInertia = 1 / mc.Inertia
```

### Physics Calculations
The engine now supports proper rotational physics:

1. **Angular Acceleration**: `α = τ / I` (torque divided by inertia)
2. **Angular Velocity**: `ω = ω₀ + α * dt` (angular velocity plus acceleration times time)
3. **Rotation**: `θ = θ₀ + ω * dt` (rotation plus angular velocity times time)

## Demo Objects

The engine now includes three types of test objects:

1. **Test Square**: Rotates at 2.0 rad/s with AABB collision
2. **Test Circle**: Rotates at -1.5 rad/s with circle collision  
3. **Rotating Object**: Starts with no rotation but has continuous torque applied

## System Order

The systems are executed in this order:
1. `UpdateCamera` - Updates camera position
2. `UpdateCollisions` - Handles collision detection
3. `UpdateVelocity` - Updates linear velocity
4. `ApplyContinuousTorque` - Applies continuous torque (demo)
5. `UpdateTorque` - Converts torque to angular velocity changes
6. `UpdateAngularVelocity` - Updates rotation based on angular velocity

This creates a complete rotational physics system that integrates with the existing linear physics. 
# Collision System Fixes

This document describes the fixes implemented to resolve the issue where objects would sometimes pull towards each other during collisions.

## Problem Analysis

The pulling issue was caused by several factors:

1. **Numerical Instability**: Division by very small numbers causing extreme impulse values
2. **Incorrect Collision Point Calculation**: Using center points instead of proper collision points
3. **Missing Impulse Clamping**: No limits on impulse values leading to unrealistic forces
4. **Poor Positional Correction**: Overly aggressive correction causing objects to stick together

## Fixes Implemented

### 1. Improved Numerical Stability

#### Before:
```go
if velAlongNormal <= 0 {
    // No threshold check
}
```

#### After:
```go
if velAlongNormal <= -0.001 { // Small threshold to prevent numerical issues
    // Only resolve if objects are moving towards each other
}
```

### 2. Better Collision Point Calculation

#### For Circles:
```go
// Calculate collision point (weighted by radius)
totalRadius := crcl1.Radius + crcl2.Radius
ratio1 := crcl2.Radius / totalRadius
ratio2 := crcl1.Radius / totalRadius
collisionPoint := Vec2.Vec2{
    X: tr1.Pos.X*ratio2 + tr2.Pos.X*ratio1,
    Y: tr1.Pos.Y*ratio2 + tr2.Pos.Y*ratio1,
}
```

### 3. Impulse Clamping

#### Linear Impulse Clamping:
```go
// Clamp impulse to prevent extreme values
maxImpulse := 1000.0 // Adjust based on your physics scale
if math.Abs(j) > maxImpulse {
    if j > 0 {
        j = maxImpulse
    } else {
        j = -maxImpulse
    }
}
```

#### Friction Impulse Clamping:
```go
// Clamp friction impulse
maxFrictionImpulse := 500.0 // Adjust based on your physics scale
if frictionImpulse.Magnitude() > maxFrictionImpulse {
    frictionImpulse.Normalize()
    frictionImpulse = frictionImpulse.Mult(maxFrictionImpulse)
}
```

### 4. Improved Positional Correction

#### Before:
```go
func PositionalCorrection(e1, e2 *donburi.Entry, n Vec2.Vec2, penetration_depth, percent float64) {
    if penetration_depth < 0 {
        return
    }
    // No additional checks
}
```

#### After:
```go
func ImprovedPositionalCorrection(e1, e2 *donburi.Entry, n Vec2.Vec2, penetration_depth, percent float64) {
    if penetration_depth < 0.001 { // Small threshold to prevent unnecessary corrections
        return
    }
    
    // Calculate correction with improved stability
    totalInverseMass := m1.InverseMass + m2.InverseMass
    if totalInverseMass < 0.001 {
        return
    }
    
    correction := n.Mult(percent * penetration_depth / totalInverseMass)
    
    // Clamp correction to prevent extreme values
    maxCorrection := 10.0 // Adjust based on your physics scale
    if correction.Magnitude() > maxCorrection {
        correction.Normalize()
        correction = correction.Mult(maxCorrection)
    }
}
```

### 5. Division by Zero Protection

#### Impulse Calculation:
```go
// Prevent division by zero and clamp impulse
if denominator > 0.001 {
    j /= denominator
    // Apply impulse...
}
```

#### Friction Calculation:
```go
if denominator > 0.001 {
    jt /= denominator
    // Apply friction...
}
```

## Key Improvements

### 1. **Velocity Threshold**
- Only resolve collisions when objects are moving towards each other
- Prevents resolution of separating objects

### 2. **Impulse Limits**
- Clamps both linear and angular impulses
- Prevents unrealistic forces from numerical errors

### 3. **Positional Correction Limits**
- Clamps positional corrections
- Prevents objects from jumping large distances

### 4. **Better Collision Points**
- Weighted collision points for circles
- More accurate impulse application

### 5. **Numerical Stability**
- Division by zero protection
- Small thresholds to prevent floating-point issues

## System Configuration

The improved collision system is now used in the scene:

```go
ms.ecs.AddSystem(systems.UpdateImprovedCollisions) // Instead of UpdateRotatedCollisions
```

## Performance Impact

- **Minimal**: The additional checks add negligible overhead
- **Improved Stability**: Prevents physics simulation from breaking down
- **Better Accuracy**: More realistic collision responses

## Testing

The fixes can be tested by:

1. **Running the demo**: Objects should no longer pull towards each other
2. **High-speed collisions**: Should remain stable
3. **Stacked objects**: Should not stick together
4. **Rotating collisions**: Should work smoothly with rotation

## Configuration

You can adjust the clamping values based on your physics scale:

```go
maxImpulse := 1000.0        // Linear impulse limit
maxFrictionImpulse := 500.0  // Friction impulse limit
maxCorrection := 10.0        // Positional correction limit
```

These values should be tuned based on your object masses and velocities. 
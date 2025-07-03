package utils

import (
	"reservio/config"
	"reservio/models"
)

// BusinessLogicValidator handles complex business rule validations
type BusinessLogicValidator struct{}

// NewBusinessLogicValidator creates a new business logic validator
func NewBusinessLogicValidator() *BusinessLogicValidator {
	return &BusinessLogicValidator{}
}

// ValidateChildOwnership checks if a child belongs to the specified parent
func (v *BusinessLogicValidator) ValidateChildOwnership(childID, parentID uint) error {
	var child models.Child
	if err := config.DB.Where("id = ? AND parent_id = ?", childID, parentID).First(&child).Error; err != nil {
		return NewValidationError(ErrChildNotFound, "Child not found or does not belong to you", map[string]interface{}{
			"child_id":  childID,
			"parent_id": parentID,
		})
	}
	return nil
}

// ValidateSlotExists checks if a slot exists
func (v *BusinessLogicValidator) ValidateSlotExists(slotID uint) error {
	var slot models.Slot
	if err := config.DB.First(&slot, slotID).Error; err != nil {
		return NewValidationError(ErrSlotNotFound, "Slot not found", map[string]interface{}{
			"slot_id": slotID,
		})
	}
	return nil
}

// ValidateReservationExists checks if a reservation exists
func (v *BusinessLogicValidator) ValidateReservationExists(reservationID uint) error {
	var reservation models.Reservation
	if err := config.DB.First(&reservation, reservationID).Error; err != nil {
		return NewValidationError(ErrReservationNotFound, "Reservation not found", map[string]interface{}{
			"reservation_id": reservationID,
		})
	}
	return nil
}

// ValidateSlotCapacity checks if a slot has available capacity
func (v *BusinessLogicValidator) ValidateSlotCapacity(slotID uint) error {
	var slot models.Slot
	if err := config.DB.First(&slot, slotID).Error; err != nil {
		return NewValidationError(ErrSlotNotFound, "Slot not found", map[string]interface{}{
			"slot_id": slotID,
		})
	}

	// Count existing reservations for this slot
	var reservationCount int64
	config.DB.Model(&models.Reservation{}).Where("slot_id = ? AND status != 'rejected'", slotID).Count(&reservationCount)

	if int(reservationCount) >= slot.Capacity {
		return NewValidationError(ErrSlotFull, "Slot is at full capacity", map[string]interface{}{
			"slot_id":          slotID,
			"current_bookings": reservationCount,
			"max_capacity":     slot.Capacity,
			"available_slots":  0,
		})
	}

	return nil
}

// ValidateDoubleBooking checks if a child is already booked for the same slot
func (v *BusinessLogicValidator) ValidateDoubleBooking(childID, slotID uint) error {
	var existingReservation models.Reservation
	if err := config.DB.Where("child_id = ? AND slot_id = ? AND status != 'rejected'", childID, slotID).First(&existingReservation).Error; err == nil {
		return NewValidationError(ErrDoubleBooking, "Child is already booked for this slot", map[string]interface{}{
			"child_id":                childID,
			"slot_id":                 slotID,
			"existing_reservation_id": existingReservation.ID,
		})
	}
	return nil
}

// ValidateReservationOwnership checks if a reservation belongs to the specified parent
func (v *BusinessLogicValidator) ValidateReservationOwnership(reservationID, parentID uint) error {
	var reservation models.Reservation
	if err := config.DB.Joins("JOIN children ON children.id = reservations.child_id").
		Where("reservations.id = ? AND children.parent_id = ?", reservationID, parentID).
		First(&reservation).Error; err != nil {
		return NewValidationError(ErrReservationNotFound, "Reservation not found or does not belong to you", map[string]interface{}{
			"reservation_id": reservationID,
			"parent_id":      parentID,
		})
	}
	return nil
}

// ValidateReservationStatusTransition checks if a status transition is valid
func (v *BusinessLogicValidator) ValidateReservationStatusTransition(currentStatus, newStatus string) error {
	validTransitions := map[string][]string{
		"pending":  {"approved", "rejected"},
		"approved": {"rejected"}, // Can still be rejected after approval
		"rejected": {},           // No further transitions from rejected
	}

	allowedTransitions, exists := validTransitions[currentStatus]
	if !exists {
		return NewValidationError(ErrInvalidStatus, "Invalid current status", map[string]interface{}{
			"current_status": currentStatus,
			"valid_statuses": []string{"pending", "approved", "rejected"},
		})
	}

	for _, allowed := range allowedTransitions {
		if newStatus == allowed {
			return nil
		}
	}

	return NewValidationError(ErrInvalidStatus, "Invalid status transition", map[string]interface{}{
		"current_status":      currentStatus,
		"new_status":          newStatus,
		"allowed_transitions": allowedTransitions,
	})
}

// ValidateUserRole checks if a user has the required role
func (v *BusinessLogicValidator) ValidateUserRole(userID uint, requiredRole string) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return NewValidationError(ErrNotFound, "User not found", map[string]interface{}{
			"user_id": userID,
		})
	}

	if user.Role != requiredRole {
		return NewValidationError(ErrForbidden, "Insufficient permissions", map[string]interface{}{
			"user_role":     user.Role,
			"required_role": requiredRole,
			"user_id":       userID,
		})
	}

	return nil
}

// GetSlotAvailability returns slot availability information
func (v *BusinessLogicValidator) GetSlotAvailability(slotID uint) (map[string]interface{}, error) {
	var slot models.Slot
	if err := config.DB.First(&slot, slotID).Error; err != nil {
		return nil, NewValidationError(ErrSlotNotFound, "Slot not found", map[string]interface{}{
			"slot_id": slotID,
		})
	}

	var reservationCount int64
	config.DB.Model(&models.Reservation{}).Where("slot_id = ? AND status != 'rejected'", slotID).Count(&reservationCount)

	availableSlots := slot.Capacity - int(reservationCount)
	if availableSlots < 0 {
		availableSlots = 0
	}

	return map[string]interface{}{
		"slot_id":         slotID,
		"total_capacity":  slot.Capacity,
		"booked_slots":    reservationCount,
		"available_slots": availableSlots,
		"is_full":         availableSlots == 0,
	}, nil
}

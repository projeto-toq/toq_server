package handlervalidators

import (
	"context"
	"log/slog"
	"time"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	crecimodel "github.com/giulio-alfieri/toq_server/internal/core/model/creci_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/converters"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CleanAndValidateProfile(ctx context.Context, in *pb.User) (user usermodel.UserInterface, err error) {
	user = usermodel.NewUser()

	if err = validateAndSetBasicInfo(user, in); err != nil {
		return
	}

	if err = validateAndSetRealtorInfo(user, in); err != nil {
		return
	}

	if err = validateAndSetContactInfo(user, in); err != nil {
		return
	}

	if err = validateAndSetAddressInfo(user, in); err != nil {
		return
	}

	user.SetPassword(converters.TrimSpaces(in.GetPassword()))

	return
}

func validateAndSetBasicInfo(user usermodel.UserInterface, in *pb.User) error {
	// user.SetFullName(converters.TrimSpaces(in.GetFullName()))
	// if err := validators.ValidateNoSpecialCharacters(user.GetFullName()); err != nil {
	// 	return err
	// }

	user.SetNickName(converters.TrimSpaces(in.GetNickName()))
	if err := validators.ValidateNoSpecialCharacters(user.GetNickName()); err != nil {
		return err
	}

	user.SetNationalID(converters.RemoveAllButDigits(in.GetNationalID()))

	date, err := time.Parse("2006-01-02", in.GetBornAT())
	if err != nil {
		slog.Error("Error converting owner bornAt to date", "error", err)
		return status.Error(codes.InvalidArgument, "Invalid date format")
	}
	user.SetBornAt(date)

	return nil
}

func validateAndSetRealtorInfo(user usermodel.UserInterface, in *pb.User) error {
	if in.GetCreciNumber() == "" {
		return nil
	}

	user.SetCreciNumber(converters.RemoveAllButDigits(in.GetCreciNumber()))
	user.SetCreciState(converters.TrimSpaces(in.GetCreciState()))

	if !crecimodel.ValidStates[user.GetCreciState()] {
		return status.Error(codes.InvalidArgument, "Invalid creci state")
	}

	date, err := time.Parse("2006-01-02", in.GetCreciValidity())
	if err != nil {
		slog.Error("Error converting owner CreciValidity to date", "error", err)
		return status.Error(codes.InvalidArgument, "Invalid date format")
	}
	user.SetCreciValidity(date)

	return nil
}

// validateAndSetContactInfo validates and sets the contact information for a user.
// It processes the phone number and email address from the input protobuf User message,
// validates them, and sets them on the user object.
//
// Parameters:
//
//	user (usermodel.UserInterface): The user object to update with contact information.
//	in (*pb.User): The input protobuf User message containing the contact information.
//
// Returns:
//
//	error: An error if the phone number or email validation fails, otherwise nil.
func validateAndSetContactInfo(user usermodel.UserInterface, in *pb.User) error {
	user.SetPhoneNumber(converters.RemoveAllButDigitsAndPlusSign(in.GetPhoneNumber()))
	if err := validators.ValidateE164(user.GetPhoneNumber()); err != nil {
		return err
	}

	user.SetEmail(converters.TrimSpaces(in.GetEmail()))
	if err := validators.ValidateEmail(user.GetEmail()); err != nil {
		return err
	}

	return nil
}

// validateAndSetAddressInfo validates and sets the address information for a user.
// It processes the input fields by removing unnecessary characters and trimming spaces
// before setting them to the user object.
//
// Parameters:
//   - user: An object implementing the UserInterface, representing the user whose address
//     information is being validated and set.
//   - in: A pointer to a pb.User object containing the address information to be validated
//     and set.
//
// Returns:
//   - error: An error if any issues occur during the validation and setting process, otherwise nil.
func validateAndSetAddressInfo(user usermodel.UserInterface, in *pb.User) error {
	user.SetZipCode(converters.RemoveAllButDigits(in.GetZipCode()))
	user.SetStreet(converters.TrimSpaces(in.GetStreet()))
	user.SetNumber(converters.TrimSpaces(in.GetNumber()))
	user.SetComplement(converters.TrimSpaces(in.GetComplement()))
	user.SetNeighborhood(converters.TrimSpaces(in.GetNeighborhood()))
	user.SetCity(converters.TrimSpaces(in.GetCity()))
	user.SetState(converters.TrimSpaces(in.GetState()))

	return nil
}

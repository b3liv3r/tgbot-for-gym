package telegram

const msgHelp = `Available commands:
/help - Show available commands
/start - Start bot message
/u_create <first_name> <last_name> <phone> <email> - Create a user profile
/u_profile - Show user profile
/u_update <first_name> <last_name> <phone> <email> - Update user profile, u can use symbol "-" for skip some field.
/u_change_subscription <subscription_level> - Change user subscription level
/u_change_current_gym <gym_id> - Change user's current gym
/w_deposit <amount> - Deposit money to wallet
/w_balance - Show current wallet balance
/w_history <start_date> <end_date> - Show wallet transaction history. Date format: 02.01.2006_15:04
/s_extend <months_count> - Extend subscription
/s_data - Show subscription data
/g_list - Show list of available gyms
/g_schedules <gym_id> - Show gym schedule
/t_list_for_gym <gym_id> - Show list of trainers for a gym
/t_available_booking_list <trainer_id> - Show available booking slots for a trainer
/t_current_booking_list - Show current bookings of a trainer
/t_booking <slot_id> - Book a slot with a trainer
/t_unbooking <booking_id> - Cancel a slot booking with a trainer
/support - support contacts
`

const msgStart = "Welcome to the GymService Telegram bot! 🏋️‍♂️\n\n" +
	"This bot was developed specifically for the course diploma project. " +
	"To see the list of available commands, use the /help command."

const (
	msgUnknownCommand                = "Unknown command 🤔. Use /help to see all available commands."
	msgUserNotFound                  = "Not found profile for this telegram account. Create it with command /u_create"
	msgNotEnoughArguments            = "Wrong number of arguments, use command /help for see all commands"
	msgNotEnoughUpdateArguments      = "Wrong number of arguments. For update not all fields use '-' \n\n example: \n/u_update - - - slawaMarlow@true.com"
	msgUserUpdateError               = "Your not can update only first name or last name.\n\n example: \n/u_update - - - slawaMarlow@true.com \n /u_updates Edgar Akhmedov - -"
	msgUnknownError                  = "Unknown error. Please try again later."
	msgUserAlreadyExist              = "User already exists. If u need update use command /u_update"
	msgProfileResp                   = "Your profile:\nName: %v\nPhone: %v\nEmail: %v\nSubscription level: %v\nCurrent Gym: %v"
	msgUnknownGym                    = "Unknown Gym, use command /g_list for see all gyms"
	msgLowLvlSub                     = "Your subscription level is low level. Change it with command /u_change_subscription"
	msgInsufficientFunds             = "Insufficient funds. You can deposit money on wallet with command /w_deposit"
	msgUnknownTransactionType        = "Unknown transaction type. Use /help for see all commands"
	msgLvlAlreadyGained              = "You cannot change a subscription to the same level "
	msgUnknownTrainerOrEmptyBookings = "Unknown trainer or empty bookings. Try again later"
	msgCurrentBookingsEmpty          = "There are currently no active bookings"
	msgUnknownUnBooking              = "You don't have such a booking, use /t_current_booking_list to see your current booking"
	msgLateUnBooking                 = "it is not possible to cancel the reservation as there are 24 hours or less left before the slot starts"
	msgPastSlotBooking               = "unable to book past slot"
	msgUnknownBooking                = "Unknown booking, use /t_available_booking_list <trainer_id> to see available booking"
	msgSupport                       = "If u see some bugs or unknown errors, write him @EdgarAkh"
)

const (
	descChangeSubscription = "subscription changed"
	descRollbackStatus     = "rollback status :"
	descCurrentBalance     = "Current balance: %.2f"
	descNewBalance         = "New Balance: %.2f"
	descDeposit            = "Deposit in tg bot"
	descExtend             = "Extend subscription lvl %d on %d month"
	descBookTrainer        = "Payment for booking trainer"
	descUnBookTrainer      = "Refund for unBooking trainer"
)

package cfg

const (
	GeneralDebug = "general.debug"
	GeneralGuild = "general.guild"

	BotPrefix = "bot.prefix"
	BotToken  = "bot.token"

	UserEmail    = "user.email"
	UserPassword = "user.password"
	UserToken    = "user.token"

	EmailDomain   = "email.domain"
	EmailPort     = "email.port"
	EmailAddress  = "email.address"
	EmailPassword = "email.password"

	DbFile = "database.file"

	TriggersAddress  = "triggers.address"
	TriggersPassword = "triggers.password"

	RoleVerified  = "roles.verified"
	RoleModerator = "roles.moderator"
	RoleOnAir     = "roles.onAir"

	ChannelStudio      = "channels.studio"
	ChannelControlRoom = "channels.controlRoom"

	AutoplayAnnounce      = "autoplay.announce"
	AutoplayCommands      = "autoplay.commands"
	AutoplayIgnoredUsers  = "autoplay.ignoredUsers"
	AutoplaySlotUp        = "autoplay.slotUp"
	AutoplayUsersInStudio = "autoplay.usersInStudio"

	VerificationUserUnverified = "verification.userUnverified"
	VerificationEmailNotFound  = "verification.emailNotFound"
	VerificationEmailTaken     = "verification.takenEmail"
	VerificationEmailSent      = "verification.emailSent"
	VerificationInvalidEmail   = "verification.invalidEmail"
	VerificationEmailSubject   = "verification.emailSubject"
	VerificationEmailContents  = "verification.emailContentsFile"
	VerificationNotConfirmed   = "verification.confirmationFailed"
	VerificationAllowedEmails  = "verification.allowedEmailsFile"
	VerificationConfirmed      = "verification.confirmed"

	MsgSyntaxError           = "msg.syntaxError"
	MsgInvalidTime           = "msg.invalidTime"
	MsgCmdInvite             = "msg.inviteInvited"
	MsgCmdRegisterWrongRole  = "msg.registerWrongRole"
	MsgCmdRegisterNewShow    = "msg.registerNewShow"
	MsgCmdRegisterReplaced   = "msg.registerReplaced"
	MsgCmdShowNotFound       = "msg.showNoShowFound"
	MsgCmdShowFound          = "msg.showShowFound"
	MsgCmdShowsEmbedSet      = "msg.showsembedSet"
	MsgCmdShowsEmbedReplaced = "msg.showsembedReplaced"
	MsgCmdUninvite           = "msg.uninviteUninvited"
	MsgCmdUnregisterNotFound = "msg.unregisterNoShowFound"
	MsgCmdUnregisterDeleted  = "msg.unregisterShowDeleted"
	MsgCmdAddHostNotHost     = "msg.addHostNotHost"
	MsgCmdAddHostDone        = "msg.addHostDone"
	MsgCmdEmailNotInList     = "msg.checkEmailNotInList"
	MsgCmdEmailNotLinked     = "msg.checkEmailNotLinked"
	MsgCmdEmailLinked        = "msg.checkEmailLinked"

	MiscCurrentListenersUrl = "misc.listenersUrl"
)

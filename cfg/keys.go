package cfg

const (
	GeneralDebug = "general.debug"
	GeneralGuild = "general.guild"

	BotPrefix = "bot.prefix"
	BotToken  = "bot.token"

	UserEmail    = "user.email"
	UserPassword = "user.password"
	UserToken    = "user.token"

	DbFile = "database.file"

	RoleModerator = "roles.moderator"
	RoleOnAir     = "roles.onAir"

	ChannelStudio      = "channels.studio"
	ChannelControlRoom = "channels.controlRoom"

	AutoplayAnnounce     = "autoplay.announce"
	AutoplayCommands     = "autoplay.commands"
	AutoplayIgnoredUsers = "autoplay.ignoredUsers"

	MsgSyntaxError = "msg.syntaxError"
	MsgInvalidTime = "msg.invalidTime"

	MsgCmdInvite = "msg.invite.invited"

	MsgCmdRegisterNewShow  = "msg.register.newShow"
	MsgCmdRegisterReplaced = "msg.register.replaced"

	MsgCmdShowNotFound = "msg.show.showNotFound"
	MsgCmdShowFound    = "msg.show.showFound"

	MsgCmdUninvite = "msg.uninvite.uninvited"

	MsgCmdUnregisterNotFound = "msg.unregister.noShowFound"
	MsgCmdUnregisterDeleted  = "msg.unregister.showDeleted"

	TestingQuickerJobs = "testing.quickerJobs"
)

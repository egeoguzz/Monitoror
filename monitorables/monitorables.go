package monitorables

import (
	"github.com/monitoror/monitoror/monitorables/azuredevops"
	"github.com/monitoror/monitoror/monitorables/github"
	"github.com/monitoror/monitoror/monitorables/gitlab"
	"github.com/monitoror/monitoror/monitorables/http"
	"github.com/monitoror/monitoror/monitorables/jenkins"
	"github.com/monitoror/monitoror/monitorables/ping"
	"github.com/monitoror/monitoror/monitorables/pingdom"
	"github.com/monitoror/monitoror/monitorables/port"
	"github.com/monitoror/monitoror/monitorables/travisci"
	notify "github.com/monitoror/monitoror/notify"
	"github.com/monitoror/monitoror/store"
)

func RegisterMonitorables(s *store.Store) {

	// ------------ AZURE DEVOPS ------------
	s.Registry.RegisterMonitorable(azuredevops.NewMonitorable(s))
	// ------------ GITHUB ------------
	s.Registry.RegisterMonitorable(github.NewMonitorable(s))
	// ------------ GITLAB ------------
	s.Registry.RegisterMonitorable(gitlab.NewMonitorable(s))
	// ------------ HTTP ------------
	s.Registry.RegisterMonitorable(http.NewMonitorable(s))
	// ------------ JENKINS ------------
	s.Registry.RegisterMonitorable(jenkins.NewMonitorable(s))
	// ------------ PING ------------
	s.Registry.RegisterMonitorable(ping.NewMonitorable(s))
	// ------------ PINGDOM ------------
	s.Registry.RegisterMonitorable(pingdom.NewMonitorable(s))
	// ------------ PORT ------------
	s.Registry.RegisterMonitorable(port.NewMonitorable(s))
	// ------------ TRAVIS CI ------------
	s.Registry.RegisterMonitorable(travisci.NewMonitorable(s))

	if s.CoreConfig.Slack_Notification_Url != "nil" {
		if s.CoreConfig.Slack_Mention_List != "nil" {

			notify.Initialization_Notifier(s.CoreConfig.NamedConfigs["default"], s.CoreConfig.Slack_Notification_Url, s.CoreConfig.Slack_Mention_List, s.CoreConfig.Slack_Notification_Emoji, s.CoreConfig.Slack_Fault_Count, s.CoreConfig.Slack_Notify_Channel, s.CoreConfig.Slack_Server_Success_Msg)
		} else {

			notify.Initialization_Notifier(s.CoreConfig.NamedConfigs["default"], s.CoreConfig.Slack_Notification_Url, "nil", s.CoreConfig.Slack_Notification_Emoji, s.CoreConfig.Slack_Fault_Count, s.CoreConfig.Slack_Notify_Channel, s.CoreConfig.Slack_Server_Success_Msg)
		}
	}
}

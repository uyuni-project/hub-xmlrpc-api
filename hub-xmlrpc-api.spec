#
# spec file for package hub-xmlrpc-api
#
# Copyright (c) 2020 SUSE LLC
#
# All modifications and additions to the file contributed by third parties
# remain the property of their copyright owners, unless otherwise agreed
# upon. The license for this file, and modifications and additions to the
# file, is the same license as for the pristine package itself (unless the
# license for the pristine package is not an Open Source License, in which
# case the license is the MIT License). An "Open Source License" is a
# license that conforms to the Open Source Definition (Version 1.9)
# published by the Open Source Initiative.

# Please submit bugfixes or comments via https://bugs.opensuse.org/
#


%global provider        github
%global provider_tld    com
%global org             uyuni-project
%global project         hub-xmlrpc-api
%global provider_prefix %{provider}.%{provider_tld}/%{org}/%{project}
%global import_path     %{provider_prefix}

Name:           %{project}
Version:        0.1.4
Release:        0
Summary:        Xmlrpc API to manage Hub
License:        Apache-2.0
Group:          Applications/Internet
URL:            https://%{provider_prefix}
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  go >= 1.9
BuildRequires:  golang-packaging
BuildRequires:  rsyslog
BuildRoot:      %{_tmppath}/%{name}-%{version}-build

Requires:       logrotate
Requires:       rsyslog
Requires:       systemd

%description
Hub-xmlrpc-api package provide an API which allows access to Uyuni server functionality in a Server of Serves architecture (Hub)

%prep
%setup -q -n  %{project}-%{version}

%build
%goprep %{import_path}
%gobuild ...

%install
%goinstall

%gofilelist

%define _release_dir  %{_builddir}/%{project}-%{version}/release

# Service file for hub xmlrpc api
install -D -m 0644 %{_release_dir}/hub-xmlrpc-api.service %{buildroot}%{_unitdir}/hub-xmlrpc-api.service

# Add config files for hub
install -d -m 0755 %{buildroot}%{_sysconfdir}/hub
install -D -m 0644 %{_release_dir}/hub.conf %{buildroot}%{_sysconfdir}/hub
install -D -m 0644 %{_release_dir}/hub-xmlrpc-api-config.json  %{buildroot}%{_sysconfdir}/hub
install -d -m 0750 %{buildroot}%{_var}/log/hub

# Add syslog config to redirect logs to /var/log/hub/hub-xmlrpc-api.log
install -D -m 0644 %{_release_dir}/hub-logs.conf %{buildroot}%{_sysconfdir}/rsyslog.d/hub-logs.conf

#logrotate config
install -D -m 0644 %{_release_dir}/hub-xmlrpc-api %{buildroot}%{_sysconfdir}/logrotate.d/hub-xmlrpc-api

%pre
%service_add_pre hub-xmlrpc-api.service

%post
%service_add_post hub-xmlrpc-api.service
if [ $1 == 1 ];then
   echo "-----------------------"
   echo "### hub-xmlrpc-api does NOT start on installation. You can set it up to start automatically with systemd, by executing the following commands"
   echo " sudo systemctl enable hub-xmlrpc-api.service"
   echo ""
   echo "### You can start Hub XMLRPC API by executing the following command"
   echo " sudo systemctl start hub-xmlrpc-api.service"
   echo ""
   echo " Make sure '/etc/hub/hub-xmlrpc-api-config.json' is pointing to correct hub instance"
   echo " Logs can be viewed at /var/log/hub/hub-xmlrpc-api.log"

   /usr/bin/systemctl restart rsyslog.service > /dev/null 2>&1 || :
fi

%preun
%service_del_preun hub-xmlrpc-api.service

%postun
%service_del_postun hub-xmlrpc-api.service

%files

%defattr(-,root,root)
%doc README.md
%{_bindir}/hub-xmlrpc-api

%{_unitdir}/hub-xmlrpc-api.service
%dir %{_sysconfdir}/hub
%dir %{_var}/log/hub

%config(noreplace) %{_sysconfdir}/rsyslog.d/hub-logs.conf
%config(noreplace) %{_sysconfdir}/hub/hub-xmlrpc-api-config.json
%config(noreplace) %{_sysconfdir}/hub/hub.conf
%config(noreplace) %{_sysconfdir}/logrotate.d/hub-xmlrpc-api

%changelog
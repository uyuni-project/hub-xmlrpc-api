#
# spec file for package hub-xmlrpc-api
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
Release:        1
Summary:        Xmlrpc API to manage Hub
License:        Apache-2.0
Group:          Applications/Internet
URL:            https://%{provider_prefix}
Source0:        %{name}-%{version}.tar.gz
Source1:        hub-xmlrpc-api.service
Source2:        hub.conf
Source3:        hub-xmlrpc-api-config.json
Source4:        hub-logs.conf

BuildRequires:  go >= 1.9
BuildRequires:  golang-packaging
BuildRoot:      %{_tmppath}/%{name}-%{version}-build

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



# Service file for hub xmlrpc api
install -D -m 0644 %{SOURCE1} %{buildroot}%{_unitdir}/hub-xmlrpc-api.service

# Add config files for hub
install -d -m 0750 %{buildroot}%{_sysconfdir}/hub
install -D -m 0644 %{SOURCE2} %{buildroot}%{_sysconfdir}/hub
install -D -m 0644 %{SOURCE3} %{buildroot}%{_sysconfdir}/hub
install -d -m 0750 %{buildroot}%{_var}/log/hub


#add syslog config to redirect logs to /var/log/hub/hub-xmlrpc-api.log
install -d -m 755 %{buildroot}%{_sysconfdir}/rsyslog.d
install -D -m 0644 %{SOURCE4} %{buildroot}%{_sysconfdir}/rsyslog.d


%pre
%service_add_pre hub-xmlrpc-api.service

%post
%service_add_post hub-xmlrpc-api.service
if [ -x /bin/systemctl ] ; then
    echo "### NOT starting on installation, please execute the following statements to configure hub-xmlrpc-api to start automatically using systemd"
    echo " sudo systemctl daemon-reload"
    echo " sudo systemctl enable hub-xmlrpc-api.service"
    echo "### You can start Hub XMLRPC API by executing"
    echo " sudo systemctl start hub-xmlrpc-api.service"
fi
systemctl restart rsyslog.service > /dev/null 2>&1 || :


%preun
%service_del_preun hub-xmlrpc-api.service

%postun
%service_del_postun hub-xmlrpc-api.service


%files -f file.lst

%defattr(-,root,root)
%dir %{_sysconfdir}/rsyslog.d
%{_bindir}/hub-xmlrpc-api

%{_unitdir}/hub-xmlrpc-api.service
%dir %{_sysconfdir}/hub
%dir %{_var}/log/hub
%config(noreplace) /etc/rsyslog.d/hub-logs.conf
%config(noreplace) %{_sysconfdir}/hub/hub-xmlrpc-api-config.json
%config(noreplace) %{_sysconfdir}/hub/hub.conf

%changelog

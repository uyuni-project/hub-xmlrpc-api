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

%define _release_dir  %{_builddir}/%{project}-%{version}/release

# Service file for hub xmlrpc api
install -D -m 0644 %{_release_dir}/hub-xmlrpc-api.service %{buildroot}%{_unitdir}/hub-xmlrpc-api.service

# Add config files for hub
install -d -m 0750 %{buildroot}%{_sysconfdir}/hub
install -D -m 0644 %{_release_dir}/hub.conf %{buildroot}%{_sysconfdir}/hub
install -D -m 0644 %{_release_dir}/hub-xmlrpc-api-config.json  %{buildroot}%{_sysconfdir}/hub
install -d -m 0750 %{buildroot}%{_var}/log/hub


#add syslog config to redirect logs to /var/log/hub/hub-xmlrpc-api.log
install -d -m 755 %{buildroot}%{_sysconfdir}/rsyslog.d
install -D -m 0644 %{_release_dir}/hub-logs.conf %{buildroot}%{_sysconfdir}/rsyslog.d


%pre
%service_add_pre hub-xmlrpc-api.service

%post
%service_add_post hub-xmlrpc-api.service
if [ -x /bin/systemctl ] ; then
    echo "### NOT starting on installation, please execute the following statements to configure hub-xmlrpc-api to start automatically using systemd"
    echo " sudo systemctl daemon-reload"
    echo " sudo systemctl enable hub-xmlrpc-api.service"
    echo ""
    echo "### You can start Hub XMLRPC API by executing"
    echo " sudo systemctl start hub-xmlrpc-api.service"
    echo ""
    echo " Make sure 'etc/hub/hub-xmlrpc-api-config.json' is pointing to correct hub instance"
    echo ""
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
* Tue Mar 31 2020 Abid Mehmood <amehmood@suse.de> 0.1.4-1
- changes. (amehmood@suse.de)
- save failed login responses as well in session to avoid lookup errors
  (amehmood@suse.de)
- Refactor multicast service interfaces (mchiaradia@suse.com)
- Revert hardcode (mchiaradia@suse.com)
- Fix typo in multicast response format (mchiaradia@suse.com)
- Adapt parser test to new interface (mchiaradia@suse.com)
- Adapt mocks to new interface (mchiaradia@suse.com)
- Automatic commit of package [hub-xmlrpc-api] release [0.1.4-1].
  (amehmood@suse.de)
- move spec files at root level as needed by tito (amehmood@suse.de)
- Fix format in example scripts (mchiaradia@suse.com)
- Add consts for paths (mchiaradia@suse.com)
- move all the files which are need for the releasing the rpm are moved into
  release folder (amehmood@suse.de)
- updated spec file (amehmood@suse.de)
- update examples (amehmood@suse.de)
- add unit file (amehmood@suse.de)
- Create service component (mchiaradia@suse.com)
- Fix README file (mchiaradia@suse.com)
- Add logging and error checking (mchiaradia@suse.com)
- Parsers test (mchiaradia@suse.com)
- Refactor in xml parsing (mchiaradia@suse.com)
- Refactor in tests, and temporary commented out tests (mchiaradia@suse.com)
- Refactor in server request parsing (mchiaradia@suse.com)
- Fix in tests after Session refactoring (mchiaradia@suse.com)
- Test for Session (mchiaradia@suse.com)
- Refactor in session (mchiaradia@suse.com)
- Remove noisy log in multicast (mchiaradia@suse.com)
- Small refactor in the encoder (mchiaradia@suse.com)
- Fix for resolving method to call in multicast (mchiaradia@suse.com)
- Remove args.go file (mchiaradia@suse.com)
- Initial version for moving codec and parser to different packages
  (mchiaradia@suse.com)
- Encoder test initial version (mchiaradia@suse.com)
- Imprive tests for parsers (mchiaradia@suse.com)
- Add error checking in parsers (mchiaradia@suse.com)
- Fault encoding (mchiaradia@suse.com)
- Clean up client interface (mchiaradia@suse.com)
- MulticastService test (mchiaradia@suse.com)
- Initial version for managing dependencies between components
  (mchiaradia@suse.com)
- Rename DefaultService files (mchiaradia@suse.com)
- Rename UnicastService files (mchiaradia@suse.com)
- Rename MulticastService files (mchiaradia@suse.com)
- Rename HubService files (mchiaradia@suse.com)
- Fix when bootstraping the environment (mchiaradia@suse.com)
- Remove global variable 'apiSession' (mchiaradia@suse.com)
- Initial version for removing global variable 'conf' (mchiaradia@suse.com)
- test for fault method (amehmood@suse.de)
- Temporary: Export parsers for testing purpose (mchiaradia@suse.com)
- suggested changes (amehmood@suse.de)
- some refactoring (amehmood@suse.de)
- tests of unicast (amehmood@suse.de)
- Use stateless functions instead of object for parsers (mchiaradia@suse.com)
- CodecRequest tests (mchiaradia@suse.com)
- Fix in attachToServers test (mchiaradia@suse.com)
- Set max parameters count to 10000 as per the xmlrpc protocol
  (mchiaradia@suse.com)
- Initial tests for server_codec (mchiaradia@suse.com)
- pass arguments correctly in case of unicast (amehmood@suse.de)
- Fix in attach servers for relay mode (mchiaradia@suse.com)
- Fix when parsing args in attach servers (mchiaradia@suse.com)
- Fix typo in parsers test (mchiaradia@suse.com)
- few more tests (amehmood@suse.de)
- Arguments Parser tests (mchiaradia@suse.com)
- Use unicast args parser (mchiaradia@suse.com)
- Add parser for unicast call (mchiaradia@suse.com)
- initial tests for multicast (amehmood@suse.de)
- add initial tests for hub login modes (amehmood@suse.de)
- add 3 mocked test servers, one represents Hub and other normal servers
  attached to Hub (amehmood@suse.de)
- add  FaultInvalidCredntials fault (amehmood@suse.de)
- change to correct method for formatting (amehmood@suse.de)
- don't load config automatically (amehmood@suse.de)
- README: remove sumaform references for now (smoioli@suse.de)
-  added text for using sumaform (amehmood@suse.de)
- README: editorial changes (smoioli@suse.de)
- Update README.md (mchiaradia@suse.com)
- Initial version (amehmood@suse.de)
- added spec & changes file (amehmood@suse.de)
- Change port number in log when starting the server (mchiaradia@suse.com)
- change path (amehmood@suse.de)


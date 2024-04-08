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

%if 0%{?rhel} == 8
%global debug_package %{nil}
%endif

%if 0%{?rhel}
# Fix ERROR: No build ID note found in
%undefine _missing_build_ids_terminate_build
%endif

%global provider        github
%global provider_tld    com
%global org             uyuni-project
%global project         hub-xmlrpc-api
%global provider_prefix %{provider}.%{provider_tld}/%{org}/%{project}

Name:           %{project}
Version:        0.8
Release:        0
Summary:        XMLRPC API for Hub environments
License:        Apache-2.0
Group:          Applications/Internet
URL:            https://%{provider_prefix}
Source0:        %{name}-%{version}.tar.gz
Source1:        vendor.tar.gz

%if 0%{?suse_version}
%if 0%{?is_opensuse}
BuildRequires:  golang(API) = 1.18
%else
BuildRequires:  go1.18-openssl
%endif
%else
BuildRequires:  go >= 1.17
%endif
BuildRequires:  golang-packaging
BuildRequires:  rsyslog
BuildRoot:      %{_tmppath}/%{name}-%{version}-build

Requires:       logrotate
Requires:       rsyslog
Requires:       systemd

%description
The Hub XMLRPC API provides an API which allows access to Uyuni functionality in a Server of Serves architecture (Hub)

%prep
%autosetup
tar -zxf %{SOURCE1}

%build
export GOFLAGS=-mod=vendor
%goprep %{provider_prefix}
%gobuild ...

%install
%goinstall


%define _release_dir  %{_builddir}/%{project}-%{version}/release

# Service file for hub xmlrpc api
install -D -m 0644 %{_release_dir}/hub-xmlrpc-api.service %{buildroot}%{_unitdir}/hub-xmlrpc-api.service

# Add config files for hub
install -d -m 0750 %{buildroot}%{_sysconfdir}/hub
install -D -m 0644 %{_release_dir}/hub.conf %{buildroot}%{_sysconfdir}/hub
install -d -m 0750 %{buildroot}%{_var}/log/hub

# Add syslog config to redirect logs to /var/log/hub/hub-xmlrpc-api.log
install -D -m 0644 %{_release_dir}/hub-logs.conf %{buildroot}%{_sysconfdir}/rsyslog.d/hub-logs.conf

#logrotate config
install -D -m 0644 %{_release_dir}/logrotate.conf %{buildroot}%{_sysconfdir}/logrotate.d/hub-xmlrpc-api

%check
# Fix OBS debug_package execution.
rm -f %{buildroot}/usr/lib/debug/%{_bindir}/%{name}-%{version}-*.debug

%pre
%if 0%{?suse_version}
%service_add_pre %{name}.service
%endif

%post
%if 0%{?rhel}
%systemd_post %{name}.service
%else
%service_add_post %{name}.service
%endif
if [ $1 == 1 ];then
   echo "-----------------------"
   echo "### The Hub XMLRPC API does NOT start on installation. You can set it up to start automatically with systemd, by executing the following commands"
   echo " sudo systemctl enable hub-xmlrpc-api.service"
   echo ""
   echo "### You can start Hub XMLRPC API by executing the following command"
   echo " sudo systemctl start hub-xmlrpc-api.service"
   echo ""
   echo " Make sure that '/etc/hub/hub.conf' is pointing to the correct Hub host"
   echo " Logs can be viewed at /var/log/hub/hub-xmlrpc-api.log"

   %if 0%{?rhel}
   %systemd_postun rsyslog.service
   %else
   %service_del_postun rsyslog.service
   %endif
fi

%preun
%if 0%{?rhel}
%systemd_preun %{name}.service
%else
%service_del_preun %{name}.service
%endif

%postun
%if 0%{?rhel}
%systemd_postun %{name}.service
%else
%service_del_postun %{name}.service
%endif

%files

%defattr(-,root,root)
%doc README.md
%{_bindir}/hub-xmlrpc-api

%{_unitdir}/hub-xmlrpc-api.service
%dir %{_sysconfdir}/hub
%dir %{_var}/log/hub

%config(noreplace) %{_sysconfdir}/rsyslog.d/hub-logs.conf
%config(noreplace) %{_sysconfdir}/hub/hub.conf
%config(noreplace) %{_sysconfdir}/logrotate.d/hub-xmlrpc-api

%changelog

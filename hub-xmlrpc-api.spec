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
Version:        0.1.2
Release:        1
Summary:        Xmlrpc API to manage Hub
License:        Apache-2.0
Group:          Hub
Url:            https://%{provider_prefix}
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


%files -f file.lst
%defattr(-,root,root)
%doc README.md
%{_bindir}/hub-xmlrpc-api



%changelog
* Fri Feb 07 2020 Abid Mehmood <amehmood@suse.de> 0.1.2-1
- Fix paths (amehmood@suse.de)

* Fri Feb 07 2020 Abid Mehmood <amehmood@suse.de> 0.1.1-1
- new package built with tito


* Fri Feb 07 2020 Abid Mehmood <amehmood@suse.de> 0.1.1-1
- new package built with tito

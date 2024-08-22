Name:           jr
Version:        0.3.9
Release:        1%{?dist}
Summary:        JR: streaming quality random data from the command line

License:        MIT
URL:            https://jrnd.io/
Source0:        https://github.com/ugol/%{name}/archive/refs/tags/v%{version}.tar.gz

BuildRequires:  golang >= 1.22.0
BuildRequires:  make

%description
JR is a CLI program that helps you to stream quality random data for your applications.

%global debug_package %{nil}

%prep
%setup -q

%build
make all %{?_smp_mflags}

%install
mkdir -p %{buildroot}/usr/bin
install -m 0755 %{_builddir}/%{name}-%{version}/build/jr %{buildroot}/usr/bin/jr

# Copy templates section
mkdir -p %{buildroot}%{_datadir}/jr/
cp -rf %{_builddir}/%{name}-%{version}/templates %{buildroot}%{_datadir}/jr/

# Copy config section
cp -rf %{_builddir}/%{name}-%{version}/config %{buildroot}%{_datadir}/jr/
  
%files
%license LICENSE
%{_bindir}/%{name}
%{_datadir}/jr/

%changelog
* Thu Aug 22 2024 Gianni Salinetti <gbsalinetti@gmail.com> - v0.3.9
- v0.3.9 release, includes jr default config files in /usr/share/jr
* Fri Aug 16 2024 Gianni Salinetti <gbsalinetti@gmail.com> - v0.3.8
- First jr package, templates included in /usr/share/jr/templates but still not seen by the program


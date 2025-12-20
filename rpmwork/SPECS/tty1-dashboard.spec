Name:           tty1-dashboard
Version:        0.1.0
Release:        1%{?dist}
Summary:        TTY1 dashboard

License:        Proprietary
Source0:        tty1-dashboard

BuildArch:      x86_64
BuildRequires:  systemd-rpm-macros
%global debug_package %{nil}
Requires(post): systemd
Requires(preun): systemd
Requires(postun): systemd

%description
TTY1 dashboard that replaces getty on tty1.

%prep
%setup -q -c -T
cp %{SOURCE0} .

%build
# Static binary, no build needed

%install
mkdir -p %{buildroot}/usr/local/bin
mkdir -p %{buildroot}%{_unitdir}

install -m 755 tty1-dashboard %{buildroot}/usr/local/bin/tty1-dashboard

cat > %{buildroot}%{_unitdir}/tty1-dashboard.service << 'EOF'
[Unit]
Description=TTY1 Dashboard
After=systemd-user-sessions.service
Conflicts=getty@tty1.service

[Service]
ExecStart=/usr/local/bin/tty1-dashboard
Restart=always
RestartSec=1
StandardInput=tty
StandardOutput=tty
TTYPath=/dev/tty1
TTYReset=yes
TTYVHangup=yes
TTYVTDisallocate=yes

[Install]
WantedBy=multi-user.target
EOF

%files
/usr/local/bin/tty1-dashboard
%{_unitdir}/tty1-dashboard.service

%post
%systemd_post tty1-dashboard.service
systemctl enable tty1-dashboard.service >/dev/null 2>&1 || :
systemctl start tty1-dashboard.service >/dev/null 2>&1 || :
systemctl disable --now getty@tty1.service >/dev/null 2>&1 || :
systemctl enable --now getty@tty2.service getty@tty3.service getty@tty4.service getty@tty5.service >/dev/null 2>&1 || :

%preun
%systemd_preun tty1-dashboard.service

%postun
%systemd_postun_with_restart tty1-dashboard.service

%changelog
* Sat Dec 20 2025 System Administrator <admin@example.com> - 0.1.0-1
- Initial package

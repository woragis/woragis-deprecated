#!/bin/bash

# Completely suppress all Gdk and GTK warnings
export GDK_SILENCE_GTK_WARNINGS=1
export GTK_SILENCE_GTK_WARNINGS=1
export G_MESSAGES_DEBUG=none
export GTK_DEBUG=no-css-cache
export GDK_DEBUG=no-css-cache

# Run Flutter and filter out all the annoying warnings
exec 2> >(grep -v "Gdk-WARNING" | grep -v "Error converting selection" | grep -v "Gtk-WARNING" | grep -v "GLib-GObject-WARNING" >&2)

flutter run "$@"

#!/bin/bash

# Suppress Gdk and GTK warnings
export GDK_SILENCE_GTK_WARNINGS=1
export GTK_SILENCE_GTK_WARNINGS=1

# Run Flutter with clean output
flutter run "$@" 2>&1 | grep -v "Gdk-WARNING" | grep -v "Error converting selection"

#!/usr/bin/env python3

import os
import re

# List of datasources to fix
datasources = [
    {
        'path': 'lib/features/blog/data/datasources/blog_local_datasource.dart',
        'class_name': 'BlogLocalDataSourceImpl',
        'injection_name': 'BlogLocalDataSource'
    },
    {
        'path': 'lib/features/about/data/datasources/about_local_datasource.dart', 
        'class_name': 'AboutLocalDataSourceImpl',
        'injection_name': 'AboutLocalDataSource'
    },
    {
        'path': 'lib/features/money/data/datasources/money_local_datasource.dart',
        'class_name': 'MoneyLocalDataSourceImpl', 
        'injection_name': 'MoneyLocalDataSource'
    },
    {
        'path': 'lib/features/testimonials/data/datasources/testimonials_local_datasource.dart',
        'class_name': 'TestimonialsLocalDataSourceImpl',
        'injection_name': 'TestimonialsLocalDataSource'
    },
    {
        'path': 'lib/features/settings/data/datasources/settings_local_datasource.dart',
        'class_name': 'SettingsLocalDataSourceImpl',
        'injection_name': 'SettingsLocalDataSource'
    }
]

def fix_datasource(file_path, class_name):
    """Fix a single datasource file"""
    full_path = f'/home/woragis/Projects/woragis/mobile/{file_path}'
    
    if not os.path.exists(full_path):
        print(f"⚠️  File not found: {full_path}")
        return False
        
    with open(full_path, 'r') as f:
        content = f.read()
    
    # Add AuthStore import if not present
    if 'import \'../../../../core/stores/auth_store.dart\';' not in content:
        # Find the last import and add after it
        import_pattern = r"(import '[^']+';)"
        imports = re.findall(import_pattern, content)
        if imports:
            last_import = imports[-1]
            auth_import = "import '../../../../core/stores/auth_store.dart';"
            content = content.replace(last_import, f"{last_import}\n{auth_import}")
    
    # Add AuthStoreBloc to constructor
    constructor_pattern = rf'class {class_name}[^{{]*\{{[^}}]*\n(.*?)(final.*?=.*?;)\n'
    match = re.search(constructor_pattern, content, re.DOTALL)
    
    if match and 'AuthStoreBloc' not in content:
        # Add AuthStoreBloc field and constructor parameter
        old_constructor = match.group(0)
        new_constructor = old_constructor.replace(
            match.group(1),
            f"{match.group(1)}  final AuthStoreBloc? _authStore;\n  \n  {class_name}({{AuthStoreBloc? authStore}}) : _authStore = authStore;\n\n  /// Ensure required NOT NULL fields are always present\n  Map<String, dynamic> _ensureRequiredFields(Map<String, dynamic> data) {{\n    final safeData = Map<String, dynamic>.from(data);\n    \n    // ✅ ENSURE user_id is always present (NOT NULL constraint)\n    if (!safeData.containsKey('user_id') || safeData['user_id'] == null || safeData['user_id'] == '') {{\n      // Try to get user ID from auth store\n      final currentUserId = _authStore?.state.user?.id;\n      if (currentUserId != null && currentUserId.isNotEmpty) {{\n        safeData['user_id'] = currentUserId;\n      }} else {{\n        // Fallback to a default user ID if not authenticated\n        safeData['user_id'] = 'system';\n      }}\n    }}\n    \n    return safeData;\n  }}\n\n"
        )
        content = content.replace(old_constructor, new_constructor)
    
    with open(full_path, 'w') as f:
        f.write(content)
    
    print(f"✅ Fixed: {file_path}")
    return True

def fix_injection_container():
    """Fix the injection container"""
    injection_path = '/home/woragis/Projects/woragis/mobile/lib/core/injection/injection_container.dart'
    
    with open(injection_path, 'r') as f:
        content = f.read()
    
    # Fix each datasource registration
    for ds in datasources:
        pattern = rf'sl\.registerLazySingleton<{ds["injection_name"]}>\(\s*\(\) => {ds["class_name"]}\(\)\s*\);'
        replacement = f'sl.registerLazySingleton<{ds["injection_name"]}>(\n    () => {ds["class_name"]}(authStore: sl<AuthStoreBloc>()),\n  );'
        
        if re.search(pattern, content, re.DOTALL):
            content = re.sub(pattern, replacement, content, flags=re.DOTALL)
            print(f"✅ Fixed injection: {ds['injection_name']}")
    
    with open(injection_path, 'w') as f:
        f.write(content)

if __name__ == '__main__':
    print("🔧 Fixing local datasources for user_id injection...")
    
    for ds in datasources:
        fix_datasource(ds['path'], ds['class_name'])
    
    fix_injection_container()
    
    print("🎉 All local datasources fixed!")


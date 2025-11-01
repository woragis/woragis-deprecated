import 'package:flutter/material.dart';
import '../../domain/entities/music_genre_entity.dart';

class MusicGenresPage extends StatefulWidget {
  const MusicGenresPage({super.key});

  @override
  State<MusicGenresPage> createState() => _MusicGenresPageState();
}

class _MusicGenresPageState extends State<MusicGenresPage> {
  final TextEditingController _searchController = TextEditingController();
  bool _isGridView = true;

  // TODO: Replace with actual music genres data from BLoC or API
  final List<MusicGenreEntity> _musicGenres = [
    MusicGenreEntity(
      id: '1',
      userId: 'user123',
      name: 'Electronic',
      description: 'Electronic music and EDM',
      order: 1,
      visible: true,
      createdAt: DateTime.now().subtract(const Duration(days: 30)),
      updatedAt: DateTime.now(),
    ),
    MusicGenreEntity(
      id: '2',
      userId: 'user123',
      name: 'Rock',
      description: 'Rock and alternative music',
      order: 2,
      visible: true,
      createdAt: DateTime.now().subtract(const Duration(days: 25)),
      updatedAt: DateTime.now(),
    ),
    MusicGenreEntity(
      id: '3',
      userId: 'user123',
      name: 'Hip Hop',
      description: 'Hip hop and rap music',
      order: 3,
      visible: true,
      createdAt: DateTime.now().subtract(const Duration(days: 20)),
      updatedAt: DateTime.now(),
    ),
    MusicGenreEntity(
      id: '4',
      userId: 'user123',
      name: 'Classical',
      description: 'Classical and orchestral music',
      order: 4,
      visible: true,
      createdAt: DateTime.now().subtract(const Duration(days: 15)),
      updatedAt: DateTime.now(),
    ),
    MusicGenreEntity(
      id: '5',
      userId: 'user123',
      name: 'Jazz',
      description: 'Jazz and blues music',
      order: 5,
      visible: true,
      createdAt: DateTime.now().subtract(const Duration(days: 10)),
      updatedAt: DateTime.now(),
    ),
  ];

  List<MusicGenreEntity> get _filteredGenres {
    if (_searchController.text.isEmpty) {
      return _musicGenres;
    }
    return _musicGenres.where((genre) =>
        genre.name.toLowerCase().contains(_searchController.text.toLowerCase()) ||
        (genre.description?.toLowerCase().contains(_searchController.text.toLowerCase()) ?? false)
    ).toList();
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Music Genres'),
        actions: [
          IconButton(
            icon: Icon(_isGridView ? Icons.list : Icons.grid_view),
            onPressed: () {
              setState(() {
                _isGridView = !_isGridView;
              });
            },
          ),
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () {
              // TODO: Navigate to add music genre page
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('Add music genre functionality coming soon!')),
              );
            },
          ),
        ],
      ),
      body: Column(
        children: [
          // Search Bar
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: 'Search music genres...',
                prefixIcon: const Icon(Icons.search),
                suffixIcon: _searchController.text.isNotEmpty
                    ? IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          _searchController.clear();
                          setState(() {});
                        },
                      )
                    : null,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              onChanged: (value) => setState(() {}),
            ),
          ),

          // Music Genres List
          Expanded(
            child: _filteredGenres.isEmpty
                ? Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(
                          Icons.music_note_outlined,
                          size: 64,
                          color: Colors.grey.shade400,
                        ),
                        const SizedBox(height: 16),
                        const Text(
                          'No music genres found',
                          style: TextStyle(fontSize: 18, fontWeight: FontWeight.w500),
                        ),
                        const SizedBox(height: 8),
                        const Text(
                          'Add your first music genre',
                          style: TextStyle(color: Colors.grey),
                        ),
                        const SizedBox(height: 24),
                        ElevatedButton.icon(
                          onPressed: () {
                            // TODO: Navigate to add music genre page
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(content: Text('Add music genre functionality coming soon!')),
                            );
                          },
                          icon: const Icon(Icons.add),
                          label: const Text('Add Genre'),
                        ),
                      ],
                    ),
                  )
                : _isGridView
                    ? _buildGridView(_filteredGenres)
                    : _buildListView(_filteredGenres),
          ),
        ],
      ),
    );
  }

  Widget _buildGridView(List<MusicGenreEntity> genres) {
    return GridView.builder(
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 1.2,
        crossAxisSpacing: 16,
        mainAxisSpacing: 16,
      ),
      itemCount: genres.length,
      itemBuilder: (context, index) {
        final genre = genres[index];
        return _buildGenreCard(genre);
      },
    );
  }

  Widget _buildListView(List<MusicGenreEntity> genres) {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: genres.length,
      itemBuilder: (context, index) {
        final genre = genres[index];
        return _buildGenreListItem(genre);
      },
    );
  }

  Widget _buildGenreCard(MusicGenreEntity genre) {
    return Card(
      child: InkWell(
        onTap: () {
          // TODO: Navigate to genre detail page
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Viewing ${genre.name}')),
          );
        },
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Icon and Name
              Row(
                children: [
                  Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: _getGenreColor(genre.name).withValues(alpha: 0.1),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Icon(
                      Icons.music_note,
                      color: _getGenreColor(genre.name),
                      size: 20,
                    ),
                  ),
                  const Spacer(),
                  if (genre.visible)
                    Icon(
                      Icons.visibility,
                      size: 16,
                      color: Colors.green.shade600,
                    ),
                ],
              ),
              const SizedBox(height: 12),

              // Name
              Text(
                genre.name,
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                ),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
              const SizedBox(height: 4),

              // Description
              if (genre.description != null)
                Text(
                  genre.description!,
                  style: const TextStyle(
                    fontSize: 12,
                    color: Colors.grey,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),

              const Spacer(),

              // Order
              Text(
                'Order: ${genre.order}',
                style: const TextStyle(
                  fontSize: 10,
                  color: Colors.grey,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildGenreListItem(MusicGenreEntity genre) {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: _getGenreColor(genre.name).withValues(alpha: 0.1),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Icon(
            Icons.music_note,
            color: _getGenreColor(genre.name),
            size: 24,
          ),
        ),
        title: Text(
          genre.name,
          style: const TextStyle(fontWeight: FontWeight.w600),
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            if (genre.description != null)
              Text(
                genre.description!,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            const SizedBox(height: 4),
            Text(
              'Order: ${genre.order}',
              style: const TextStyle(fontSize: 12, color: Colors.grey),
            ),
          ],
        ),
        trailing: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            if (genre.visible)
              Icon(
                Icons.visibility,
                size: 16,
                color: Colors.green.shade600,
              ),
            const Icon(Icons.arrow_forward_ios, size: 16),
          ],
        ),
        onTap: () {
          // TODO: Navigate to genre detail page
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Viewing ${genre.name}')),
          );
        },
      ),
    );
  }

  Color _getGenreColor(String genreName) {
    // Generate a consistent color based on genre name
    final colors = [
      Colors.red,
      Colors.blue,
      Colors.green,
      Colors.orange,
      Colors.purple,
      Colors.teal,
      Colors.pink,
      Colors.indigo,
      Colors.amber,
      Colors.cyan,
    ];
    
    final index = genreName.hashCode.abs() % colors.length;
    return colors[index];
  }
}

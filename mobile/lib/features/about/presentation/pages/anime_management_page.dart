import 'package:flutter/material.dart';
import '../../domain/entities/anime_entity.dart';

class AnimeManagementPage extends StatefulWidget {
  const AnimeManagementPage({super.key});

  @override
  State<AnimeManagementPage> createState() => _AnimeManagementPageState();
}

class _AnimeManagementPageState extends State<AnimeManagementPage> {
  final TextEditingController _searchController = TextEditingController();
  AnimeStatus? _selectedStatus;
  bool _isGridView = true;

  // TODO: Replace with actual anime data from BLoC or API
  final List<AnimeEntity> _animeList = [
    AnimeEntity(
      id: '1',
      userId: 'user123',
      title: 'Attack on Titan',
      description: 'Humanity fights for survival against the Titans',
      status: AnimeStatus.completed,
      myAnimeListId: '16498',
      coverImage: 'https://example.com/aot.jpg',
      genres: ['Action', 'Drama', 'Fantasy'],
      episodes: 25,
      currentEpisode: 25,
      rating: 9.5,
      notes: 'Amazing series with great character development',
      startedAt: DateTime(2023, 1, 1),
      completedAt: DateTime(2023, 3, 15),
      order: 1,
      visible: true,
      createdAt: DateTime.now().subtract(const Duration(days: 30)),
      updatedAt: DateTime.now(),
    ),
    AnimeEntity(
      id: '2',
      userId: 'user123',
      title: 'Demon Slayer',
      description: 'A young boy becomes a demon slayer to save his sister',
      status: AnimeStatus.watching,
      myAnimeListId: '38000',
      coverImage: 'https://example.com/demon-slayer.jpg',
      genres: ['Action', 'Supernatural', 'Historical'],
      episodes: 26,
      currentEpisode: 15,
      rating: 8.8,
      notes: 'Beautiful animation and compelling story',
      startedAt: DateTime(2023, 4, 1),
      order: 2,
      visible: true,
      createdAt: DateTime.now().subtract(const Duration(days: 20)),
      updatedAt: DateTime.now(),
    ),
  ];

  List<AnimeEntity> get _filteredAnime {
    var filtered = _animeList;

    if (_searchController.text.isNotEmpty) {
      filtered = filtered.where((anime) =>
          anime.title.toLowerCase().contains(_searchController.text.toLowerCase()) ||
          (anime.description?.toLowerCase().contains(_searchController.text.toLowerCase()) ?? false)
      ).toList();
    }

    if (_selectedStatus != null) {
      filtered = filtered.where((anime) => anime.status == _selectedStatus).toList();
    }

    return filtered;
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
        title: const Text('Anime & Shows'),
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
              // TODO: Navigate to add anime page
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('Add anime functionality coming soon!')),
              );
            },
          ),
        ],
      ),
      body: Column(
        children: [
          // Search and Filters
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              children: [
                // Search Bar
                TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    hintText: 'Search anime...',
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
                const SizedBox(height: 12),

                // Filter Chips
                SingleChildScrollView(
                  scrollDirection: Axis.horizontal,
                  child: Row(
                    children: [
                      // Status Filter
                      FilterChip(
                        label: Text(_selectedStatus?.name.toUpperCase() ?? 'All Status'),
                        selected: _selectedStatus != null,
                        onSelected: (selected) {
                          setState(() {
                            _selectedStatus = selected ? AnimeStatus.watching : null;
                          });
                        },
                      ),
                      const SizedBox(width: 8),

                      // Clear Filters
                      if (_selectedStatus != null || _searchController.text.isNotEmpty)
                        ActionChip(
                          label: const Text('Clear'),
                          onPressed: () {
                            setState(() {
                              _searchController.clear();
                              _selectedStatus = null;
                            });
                          },
                          backgroundColor: Colors.red.shade50,
                          labelStyle: TextStyle(color: Colors.red.shade700),
                        ),
                    ],
                  ),
                ),
              ],
            ),
          ),

          // Anime List
          Expanded(
            child: _filteredAnime.isEmpty
                ? Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(
                          Icons.movie_outlined,
                          size: 64,
                          color: Colors.grey.shade400,
                        ),
                        const SizedBox(height: 16),
                        const Text(
                          'No anime found',
                          style: TextStyle(fontSize: 18, fontWeight: FontWeight.w500),
                        ),
                        const SizedBox(height: 8),
                        const Text(
                          'Add your first anime or show',
                          style: TextStyle(color: Colors.grey),
                        ),
                        const SizedBox(height: 24),
                        ElevatedButton.icon(
                          onPressed: () {
                            // TODO: Navigate to add anime page
                            ScaffoldMessenger.of(context).showSnackBar(
                              const SnackBar(content: Text('Add anime functionality coming soon!')),
                            );
                          },
                          icon: const Icon(Icons.add),
                          label: const Text('Add Anime'),
                        ),
                      ],
                    ),
                  )
                : _isGridView
                    ? _buildGridView(_filteredAnime)
                    : _buildListView(_filteredAnime),
          ),
        ],
      ),
    );
  }

  Widget _buildGridView(List<AnimeEntity> animeList) {
    return GridView.builder(
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 0.8,
        crossAxisSpacing: 16,
        mainAxisSpacing: 16,
      ),
      itemCount: animeList.length,
      itemBuilder: (context, index) {
        final anime = animeList[index];
        return _buildAnimeCard(anime);
      },
    );
  }

  Widget _buildListView(List<AnimeEntity> animeList) {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: animeList.length,
      itemBuilder: (context, index) {
        final anime = animeList[index];
        return _buildAnimeListItem(anime);
      },
    );
  }

  Widget _buildAnimeCard(AnimeEntity anime) {
    return Card(
      child: InkWell(
        onTap: () {
          // TODO: Navigate to anime detail page
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Viewing ${anime.title}')),
          );
        },
        borderRadius: BorderRadius.circular(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Cover Image
            Expanded(
              flex: 3,
              child: Container(
                width: double.infinity,
                decoration: BoxDecoration(
                  color: Colors.grey.shade200,
                  borderRadius: const BorderRadius.vertical(top: Radius.circular(12)),
                ),
                child: anime.coverImage != null
                    ? ClipRRect(
                        borderRadius: const BorderRadius.vertical(top: Radius.circular(12)),
                        child: Image.network(
                          anime.coverImage!,
                          fit: BoxFit.cover,
                          errorBuilder: (context, error, stackTrace) {
                            return const Icon(Icons.movie, size: 40);
                          },
                        ),
                      )
                    : const Icon(Icons.movie, size: 40),
              ),
            ),

            // Content
            Expanded(
              flex: 2,
              child: Padding(
                padding: const EdgeInsets.all(12),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Title
                    Text(
                      anime.title,
                      style: const TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),

                    // Status
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                      decoration: BoxDecoration(
                        color: _getAnimeStatusColor(anime.status).withValues(alpha: 0.1),
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Text(
                        anime.status.name.toUpperCase(),
                        style: TextStyle(
                          fontSize: 8,
                          fontWeight: FontWeight.bold,
                          color: _getAnimeStatusColor(anime.status),
                        ),
                      ),
                    ),
                    const Spacer(),

                    // Progress
                    if (anime.episodes != null && anime.currentEpisode != null)
                      LinearProgressIndicator(
                        value: anime.currentEpisode! / anime.episodes!,
                        backgroundColor: Colors.grey.shade300,
                        valueColor: AlwaysStoppedAnimation<Color>(
                          _getAnimeStatusColor(anime.status),
                        ),
                      ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAnimeListItem(AnimeEntity anime) {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: Colors.grey.shade200,
            borderRadius: BorderRadius.circular(8),
          ),
          child: anime.coverImage != null
              ? ClipRRect(
                  borderRadius: BorderRadius.circular(8),
                  child: Image.network(
                    anime.coverImage!,
                    fit: BoxFit.cover,
                    errorBuilder: (context, error, stackTrace) {
                      return const Icon(Icons.movie, size: 24);
                    },
                  ),
                )
              : const Icon(Icons.movie, size: 24),
        ),
        title: Text(
          anime.title,
          style: const TextStyle(fontWeight: FontWeight.w600),
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            if (anime.description != null)
              Text(
                anime.description!,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            const SizedBox(height: 4),
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color: _getAnimeStatusColor(anime.status).withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    anime.status.name.toUpperCase(),
                    style: TextStyle(
                      fontSize: 10,
                      fontWeight: FontWeight.bold,
                      color: _getAnimeStatusColor(anime.status),
                    ),
                  ),
                ),
                if (anime.rating != null) ...[
                  const SizedBox(width: 8),
                  Row(
                    children: [
                      const Icon(Icons.star, size: 12, color: Colors.amber),
                      const SizedBox(width: 2),
                      Text(
                        anime.rating!.toStringAsFixed(1),
                        style: const TextStyle(fontSize: 10),
                      ),
                    ],
                  ),
                ],
              ],
            ),
          ],
        ),
        trailing: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            if (anime.visible)
              Icon(
                Icons.visibility,
                size: 16,
                color: Colors.green.shade600,
              ),
            const Icon(Icons.arrow_forward_ios, size: 16),
          ],
        ),
        onTap: () {
          // TODO: Navigate to anime detail page
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Viewing ${anime.title}')),
          );
        },
      ),
    );
  }

  Color _getAnimeStatusColor(AnimeStatus status) {
    switch (status) {
      case AnimeStatus.watching:
        return Colors.blue;
      case AnimeStatus.completed:
        return Colors.green;
      case AnimeStatus.onHold:
        return Colors.orange;
      case AnimeStatus.dropped:
        return Colors.red;
      case AnimeStatus.planToWatch:
        return Colors.purple;
    }
  }
}

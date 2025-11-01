import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../../domain/entities/blog_post_entity.dart';
import '../bloc/blog_bloc.dart';

class BlogPostsListPage extends StatefulWidget {
  const BlogPostsListPage({super.key});

  @override
  State<BlogPostsListPage> createState() => _BlogPostsListPageState();
}

class _BlogPostsListPageState extends State<BlogPostsListPage> {
  final _searchController = TextEditingController();
  bool _showPublishedOnly = true;
  bool _showFeaturedOnly = false;
  bool _showPublicOnly = false;
  bool _isGridView = false;
  final List<String> _selectedTags = [];

  @override
  void initState() {
    super.initState();
    // Load blog posts when page initializes
    _loadBlogPosts();
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _loadBlogPosts() {
    final searchQuery = _searchController.text.trim();
    context.read<BlogBloc>().add(GetBlogPostsRequested(
      search: searchQuery.isEmpty ? null : searchQuery,
      published: _showPublishedOnly,
      featured: _showFeaturedOnly,
      public: _showPublicOnly,
    ));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Blog Posts'),
        backgroundColor: Colors.green.shade600,
        foregroundColor: Colors.white,
        elevation: 0,
        actions: [
          IconButton(
            onPressed: () => _showFilterDialog(context),
            icon: const Icon(Icons.filter_list),
          ),
          IconButton(
            onPressed: () {
              setState(() {
                _isGridView = !_isGridView;
              });
            },
            icon: Icon(_isGridView ? Icons.list : Icons.grid_view),
          ),
        ],
      ),
      body: Column(
        children: [
          // Search bar
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: 'Search blog posts...',
                prefixIcon: const Icon(Icons.search),
                suffixIcon: _searchController.text.isNotEmpty
                    ? IconButton(
                        onPressed: () {
                          _searchController.clear();
                          _performSearch();
                        },
                        icon: const Icon(Icons.clear),
                      )
                    : null,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              onChanged: (value) => _performSearch(),
            ),
          ),

          // Filter chips
          if (_showPublishedOnly || _showFeaturedOnly || _showPublicOnly || _selectedTags.isNotEmpty)
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              height: 50,
              child: ListView(
                scrollDirection: Axis.horizontal,
                children: [
                  if (_showPublishedOnly)
                    Chip(
                      label: const Text('Published'),
                      onDeleted: () => _removeFilter('published'),
                      deleteIcon: const Icon(Icons.close, size: 18),
                    ),
                  if (_showFeaturedOnly)
                    Chip(
                      label: const Text('Featured'),
                      onDeleted: () => _removeFilter('featured'),
                      deleteIcon: const Icon(Icons.close, size: 18),
                    ),
                  if (_showPublicOnly)
                    Chip(
                      label: const Text('Public'),
                      onDeleted: () => _removeFilter('public'),
                      deleteIcon: const Icon(Icons.close, size: 18),
                    ),
                  ..._selectedTags.map((tag) => Chip(
                    label: Text(tag),
                    onDeleted: () => _removeTag(tag),
                    deleteIcon: const Icon(Icons.close, size: 18),
                  )),
                ],
              ),
            ),

          // Blog posts list
          Expanded(
            child: BlocBuilder<BlogBloc, BlogState>(
              builder: (context, state) {
                if (state is BlogLoading) {
                  return const Center(child: CircularProgressIndicator());
                } else if (state is BlogError) {
                  return Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(
                          Icons.error_outline,
                          size: 64,
                          color: Colors.red.shade300,
                        ),
                        const SizedBox(height: 16),
                        Text(
                          'Error loading blog posts',
                          style: TextStyle(
                            fontSize: 18,
                            color: Colors.red.shade700,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          state.message,
                          style: TextStyle(color: Colors.red.shade600),
                          textAlign: TextAlign.center,
                        ),
                        const SizedBox(height: 16),
                        ElevatedButton(
                          onPressed: () {
                            context.read<BlogBloc>().add(GetBlogPostsRequested());
                          },
                          child: const Text('Retry'),
                        ),
                      ],
                    ),
                  );
                } else if (state is BlogPostsLoaded) {
                  if (state.posts.isEmpty) {
                    return Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(
                            Icons.article_outlined,
                            size: 64,
                            color: Colors.grey.shade400,
                          ),
                          const SizedBox(height: 16),
                          const Text(
                            'No blog posts found',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Text(
                            'Create your first blog post to share your thoughts!',
                            style: TextStyle(
                              color: Colors.grey.shade600,
                            ),
                            textAlign: TextAlign.center,
                          ),
                        ],
                      ),
                    );
                  }

                  return RefreshIndicator(
                    onRefresh: () async {
                      context.read<BlogBloc>().add(GetBlogPostsRequested());
                    },
                    child: _isGridView
                        ? _buildGridView(state.posts)
                        : _buildListView(state.posts),
                  );
                }

                return const SizedBox.shrink();
              },
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/blog/create'),
        backgroundColor: Colors.green.shade600,
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  Widget _buildGridView(List<BlogPostEntity> posts) {
    return GridView.builder(
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 0.8,
        crossAxisSpacing: 16,
        mainAxisSpacing: 16,
      ),
      itemCount: posts.length,
      itemBuilder: (context, index) {
        final post = posts[index];
        return _buildPostCard(context, post);
      },
    );
  }

  Widget _buildListView(List<BlogPostEntity> posts) {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: posts.length,
      itemBuilder: (context, index) {
        final post = posts[index];
        return _buildPostListItem(context, post);
      },
    );
  }

  Widget _buildPostCard(BuildContext context, BlogPostEntity post) {
    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: () async {
          await context.push('/blog/${post.id}');
          // Reload list when returning from detail page
          _loadBlogPosts();
        },
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Featured image
            Expanded(
              flex: 3,
              child: Stack(
                children: [
                  Container(
                    width: double.infinity,
                    decoration: BoxDecoration(
                      image: post.featuredImage != null
                          ? DecorationImage(
                              image: NetworkImage(post.featuredImage!),
                              fit: BoxFit.cover,
                            )
                          : null,
                      color: post.featuredImage == null ? Colors.grey.shade200 : null,
                    ),
                    child: post.featuredImage == null
                        ? const Center(
                            child: Icon(
                              Icons.image,
                              size: 48,
                              color: Colors.grey,
                            ),
                          )
                        : null,
                  ),
                  if (post.featured)
                    Positioned(
                      top: 8,
                      right: 8,
                      child: Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 6,
                          vertical: 2,
                        ),
                        decoration: BoxDecoration(
                          color: Colors.amber.shade100,
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Text(
                          'Featured',
                          style: TextStyle(
                            fontSize: 10,
                            color: Colors.amber.shade800,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ),
                    ),
                  if (!post.published)
                    Positioned(
                      top: 8,
                      left: 8,
                      child: Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 6,
                          vertical: 2,
                        ),
                        decoration: BoxDecoration(
                          color: Colors.orange.shade100,
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Text(
                          'Draft',
                          style: TextStyle(
                            fontSize: 10,
                            color: Colors.orange.shade800,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ),
                    ),
                ],
              ),
            ),
            // Post info
            Expanded(
              flex: 2,
              child: Padding(
                padding: const EdgeInsets.all(12),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      post.title,
                      style: const TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.bold,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),
                    Text(
                      post.excerpt,
                      style: TextStyle(
                        fontSize: 12,
                        color: Colors.grey.shade600,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const Spacer(),
                    Row(
                      children: [
                        Icon(
                          Icons.visibility,
                          size: 12,
                          color: post.visible ? Colors.green : Colors.grey,
                        ),
                        const SizedBox(width: 4),
                        Text(
                          post.visible ? 'Visible' : 'Hidden',
                          style: TextStyle(
                            fontSize: 10,
                            color: post.visible ? Colors.green : Colors.grey,
                          ),
                        ),
                        const Spacer(),
                        if (post.readingTime != null) ...[
                          Icon(
                            Icons.schedule,
                            size: 12,
                            color: Colors.blue.shade600,
                          ),
                          const SizedBox(width: 2),
                          Text(
                            '${post.readingTime}m',
                            style: TextStyle(
                              fontSize: 10,
                              color: Colors.blue.shade600,
                            ),
                          ),
                        ],
                      ],
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

  Widget _buildPostListItem(BuildContext context, BlogPostEntity post) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: () async {
          await context.push('/blog/${post.id}');
          // Reload list when returning from detail page
          _loadBlogPosts();
        },
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              // Featured image
              ClipRRect(
                borderRadius: BorderRadius.circular(8),
                child: Container(
                  width: 80,
                  height: 80,
                  decoration: post.featuredImage != null
                      ? BoxDecoration(
                          image: DecorationImage(
                            image: NetworkImage(post.featuredImage!),
                            fit: BoxFit.cover,
                          ),
                        )
                      : BoxDecoration(
                          color: Colors.grey.shade200,
                        ),
                  child: post.featuredImage == null
                      ? const Icon(
                          Icons.image,
                          size: 32,
                          color: Colors.grey,
                        )
                      : null,
                ),
              ),
              const SizedBox(width: 16),
              // Post info
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            post.title,
                            style: const TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                        if (post.featured)
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 6,
                              vertical: 2,
                            ),
                            decoration: BoxDecoration(
                              color: Colors.amber.shade100,
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: Text(
                              'Featured',
                              style: TextStyle(
                                fontSize: 10,
                                color: Colors.amber.shade800,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ),
                        const SizedBox(width: 4),
                        if (!post.published)
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 6,
                              vertical: 2,
                            ),
                            decoration: BoxDecoration(
                              color: Colors.orange.shade100,
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: Text(
                              'Draft',
                              style: TextStyle(
                                fontSize: 10,
                                color: Colors.orange.shade800,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ),
                      ],
                    ),
                    const SizedBox(height: 4),
                    Text(
                      post.excerpt,
                      style: TextStyle(
                        color: Colors.grey.shade600,
                        fontSize: 14,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        Icon(
                          Icons.visibility,
                          size: 16,
                          color: post.visible ? Colors.green : Colors.grey,
                        ),
                        const SizedBox(width: 4),
                        Text(
                          post.visible ? 'Visible' : 'Hidden',
                          style: TextStyle(
                            fontSize: 12,
                            color: post.visible ? Colors.green : Colors.grey,
                          ),
                        ),
                        const SizedBox(width: 16),
                        Icon(
                          Icons.public,
                          size: 16,
                          color: post.public ? Colors.blue : Colors.grey,
                        ),
                        const SizedBox(width: 4),
                        Text(
                          post.public ? 'Public' : 'Private',
                          style: TextStyle(
                            fontSize: 12,
                            color: post.public ? Colors.blue : Colors.grey,
                          ),
                        ),
                        const Spacer(),
                        if (post.readingTime != null)
                          Row(
                            children: [
                              Icon(
                                Icons.schedule,
                                size: 16,
                                color: Colors.blue.shade600,
                              ),
                              const SizedBox(width: 4),
                              Text(
                                '${post.readingTime}m read',
                                style: TextStyle(
                                  fontSize: 12,
                                  color: Colors.blue.shade600,
                                ),
                              ),
                            ],
                          ),
                      ],
                    ),
                    if (post.tags != null && post.tags!.isNotEmpty) ...[
                      const SizedBox(height: 8),
                      Wrap(
                        spacing: 4,
                        runSpacing: 4,
                        children: post.tags!.take(3).map((tag) {
                          return Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 6,
                              vertical: 2,
                            ),
                            decoration: BoxDecoration(
                              color: Colors.blue.shade50,
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: Text(
                              tag.name,
                              style: TextStyle(
                                fontSize: 10,
                                color: Colors.blue.shade700,
                              ),
                            ),
                          );
                        }).toList(),
                      ),
                    ],
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _performSearch() {
    final searchQuery = _searchController.text.trim();
    context.read<BlogBloc>().add(GetBlogPostsRequested(
      search: searchQuery.isEmpty ? null : searchQuery,
      published: _showPublishedOnly,
      featured: _showFeaturedOnly,
      public: _showPublicOnly,
    ));
  }

  void _removeFilter(String filterType) {
    setState(() {
      if (filterType == 'published') {
        _showPublishedOnly = false;
      } else if (filterType == 'featured') {
        _showFeaturedOnly = false;
      } else if (filterType == 'public') {
        _showPublicOnly = false;
      }
    });
    _performSearch();
  }

  void _removeTag(String tag) {
    setState(() {
      _selectedTags.remove(tag);
    });
    _performSearch();
  }

  void _showFilterDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Filter Blog Posts'),
        content: StatefulBuilder(
          builder: (context, setDialogState) {
            return Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                CheckboxListTile(
                  title: const Text('Published only'),
                  value: _showPublishedOnly,
                  onChanged: (value) {
                    setDialogState(() {
                      _showPublishedOnly = value ?? false;
                    });
                  },
                ),
                CheckboxListTile(
                  title: const Text('Featured only'),
                  value: _showFeaturedOnly,
                  onChanged: (value) {
                    setDialogState(() {
                      _showFeaturedOnly = value ?? false;
                    });
                  },
                ),
                CheckboxListTile(
                  title: const Text('Public only'),
                  value: _showPublicOnly,
                  onChanged: (value) {
                    setDialogState(() {
                      _showPublicOnly = value ?? false;
                    });
                  },
                ),
                const Divider(),
                const Text('Tags:', style: TextStyle(fontWeight: FontWeight.bold)),
                // TODO: Add tag selection chips
                const Text('Tag filtering will be implemented soon'),
              ],
            );
          },
        ),
        actions: [
          TextButton(
            onPressed: () {
              setState(() {
                _showPublishedOnly = false;
                _showFeaturedOnly = false;
                _showPublicOnly = false;
                _selectedTags.clear();
              });
              _performSearch();
              Navigator.pop(context);
            },
            child: const Text('Clear'),
          ),
          ElevatedButton(
            onPressed: () {
              setState(() {});
              _performSearch();
              Navigator.pop(context);
            },
            child: const Text('Apply'),
          ),
        ],
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../../domain/entities/testimonial_entity.dart';
import '../bloc/testimonials_bloc.dart';

class TestimonialsListPage extends StatefulWidget {
  const TestimonialsListPage({super.key});

  @override
  State<TestimonialsListPage> createState() => _TestimonialsListPageState();
}

class _TestimonialsListPageState extends State<TestimonialsListPage> {
  final _searchController = TextEditingController();
  bool _showFeaturedOnly = false;
  bool _showPublicOnly = false;
  bool _isGridView = false;
  int? _selectedRating;

  @override
  void initState() {
    super.initState();
    // Load testimonials when page initializes
    _loadTestimonials();
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _loadTestimonials() {
    final searchQuery = _searchController.text.trim();
    context.read<TestimonialsBloc>().add(GetTestimonialsRequested(
      search: searchQuery.isEmpty ? null : searchQuery,
      featured: _showFeaturedOnly,
      public: _showPublicOnly,
      rating: _selectedRating,
    ));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Testimonials'),
        backgroundColor: Colors.purple.shade600,
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
                hintText: 'Search testimonials...',
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
          if (_showFeaturedOnly || _showPublicOnly || _selectedRating != null)
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              height: 50,
              child: ListView(
                scrollDirection: Axis.horizontal,
                children: [
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
                  if (_selectedRating != null)
                    Chip(
                      label: Text('$_selectedRating★+'),
                      onDeleted: () => _removeFilter('rating'),
                      deleteIcon: const Icon(Icons.close, size: 18),
                    ),
                ],
              ),
            ),

          // Testimonials list
          Expanded(
            child: BlocBuilder<TestimonialsBloc, TestimonialsState>(
              builder: (context, state) {
                if (state is TestimonialsLoading) {
                  return const Center(child: CircularProgressIndicator());
                } else if (state is TestimonialsError) {
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
                          'Error loading testimonials',
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
                            context.read<TestimonialsBloc>().add(GetTestimonialsRequested());
                          },
                          child: const Text('Retry'),
                        ),
                      ],
                    ),
                  );
                } else if (state is TestimonialsLoaded) {
                  if (state.testimonials.isEmpty) {
                    return Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(
                            Icons.format_quote,
                            size: 64,
                            color: Colors.grey.shade400,
                          ),
                          const SizedBox(height: 16),
                          const Text(
                            'No testimonials found',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Text(
                            'Add testimonials to showcase client feedback!',
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
                      context.read<TestimonialsBloc>().add(GetTestimonialsRequested());
                    },
                    child: _isGridView
                        ? _buildGridView(state.testimonials)
                        : _buildListView(state.testimonials),
                  );
                }

                return const SizedBox.shrink();
              },
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/testimonials/create'),
        backgroundColor: Colors.purple.shade600,
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  Widget _buildGridView(List<TestimonialEntity> testimonials) {
    return GridView.builder(
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 0.8,
        crossAxisSpacing: 16,
        mainAxisSpacing: 16,
      ),
      itemCount: testimonials.length,
      itemBuilder: (context, index) {
        final testimonial = testimonials[index];
        return _buildTestimonialCard(context, testimonial);
      },
    );
  }

  Widget _buildListView(List<TestimonialEntity> testimonials) {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: testimonials.length,
      itemBuilder: (context, index) {
        final testimonial = testimonials[index];
        return _buildTestimonialListItem(context, testimonial);
      },
    );
  }

  Widget _buildTestimonialCard(BuildContext context, TestimonialEntity testimonial) {
    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: () async {
          await context.push('/testimonials/${testimonial.id}');
          // Reload list when returning from detail page
          _loadTestimonials();
        },
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Avatar and rating
              Row(
                children: [
                  CircleAvatar(
                    radius: 20,
                    backgroundColor: Colors.purple.shade100,
                    backgroundImage: testimonial.avatar != null
                        ? NetworkImage(testimonial.avatar!)
                        : null,
                    child: testimonial.avatar == null
                        ? Text(
                            testimonial.name.isNotEmpty
                                ? testimonial.name[0].toUpperCase()
                                : 'T',
                            style: TextStyle(
                              color: Colors.purple.shade700,
                              fontWeight: FontWeight.bold,
                            ),
                          )
                        : null,
                  ),
                  const Spacer(),
                  _buildRatingStars(testimonial.rating, size: 14),
                  if (testimonial.featured)
                    Container(
                      margin: const EdgeInsets.only(left: 8),
                      padding: const EdgeInsets.symmetric(
                        horizontal: 4,
                        vertical: 2,
                      ),
                      decoration: BoxDecoration(
                        color: Colors.amber.shade100,
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Text(
                        'Featured',
                        style: TextStyle(
                          fontSize: 8,
                          color: Colors.amber.shade800,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ),
                ],
              ),
              const SizedBox(height: 12),
              // Content
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      testimonial.content,
                      style: const TextStyle(
                        fontSize: 12,
                        fontStyle: FontStyle.italic,
                      ),
                      maxLines: 4,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const Spacer(),
                    // Author info
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          testimonial.name,
                          style: const TextStyle(
                            fontSize: 11,
                            fontWeight: FontWeight.bold,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                        Text(
                          testimonial.position,
                          style: TextStyle(
                            fontSize: 10,
                            color: Colors.grey.shade600,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                        Text(
                          testimonial.company,
                          style: TextStyle(
                            fontSize: 10,
                            color: Colors.grey.shade500,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildTestimonialListItem(BuildContext context, TestimonialEntity testimonial) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: () async {
          await context.push('/testimonials/${testimonial.id}');
          // Reload list when returning from detail page
          _loadTestimonials();
        },
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              // Avatar
              CircleAvatar(
                radius: 30,
                backgroundColor: Colors.purple.shade100,
                backgroundImage: testimonial.avatar != null
                    ? NetworkImage(testimonial.avatar!)
                    : null,
                child: testimonial.avatar == null
                    ? Text(
                        testimonial.name.isNotEmpty
                            ? testimonial.name[0].toUpperCase()
                            : 'T',
                        style: TextStyle(
                          color: Colors.purple.shade700,
                          fontWeight: FontWeight.bold,
                          fontSize: 20,
                        ),
                      )
                    : null,
              ),
              const SizedBox(width: 16),
              // Content
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Rating and featured badge
                    Row(
                      children: [
                        _buildRatingStars(testimonial.rating),
                        const SizedBox(width: 8),
                        if (testimonial.featured)
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
                      ],
                    ),
                    const SizedBox(height: 8),
                    // Testimonial content
                    Text(
                      testimonial.content,
                      style: TextStyle(
                        color: Colors.grey.shade700,
                        fontSize: 14,
                        fontStyle: FontStyle.italic,
                      ),
                      maxLines: 3,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 8),
                    // Author info
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          testimonial.name,
                          style: const TextStyle(
                            fontSize: 14,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        Text(
                          testimonial.position,
                          style: TextStyle(
                            fontSize: 12,
                            color: Colors.grey.shade600,
                          ),
                        ),
                        Text(
                          testimonial.company,
                          style: TextStyle(
                            fontSize: 12,
                            color: Colors.grey.shade500,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    // Status indicators
                    Row(
                      children: [
                        Icon(
                          Icons.visibility,
                          size: 16,
                          color: testimonial.visible ? Colors.green : Colors.grey,
                        ),
                        const SizedBox(width: 4),
                        Text(
                          testimonial.visible ? 'Visible' : 'Hidden',
                          style: TextStyle(
                            fontSize: 12,
                            color: testimonial.visible ? Colors.green : Colors.grey,
                          ),
                        ),
                        const SizedBox(width: 16),
                        Icon(
                          Icons.public,
                          size: 16,
                          color: testimonial.public ? Colors.blue : Colors.grey,
                        ),
                        const SizedBox(width: 4),
                        Text(
                          testimonial.public ? 'Public' : 'Private',
                          style: TextStyle(
                            fontSize: 12,
                            color: testimonial.public ? Colors.blue : Colors.grey,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildRatingStars(int rating, {double size = 16}) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: List.generate(5, (index) {
        return Icon(
          index < rating ? Icons.star : Icons.star_border,
          size: size,
          color: Colors.amber,
        );
      }),
    );
  }

  void _performSearch() {
    final searchQuery = _searchController.text.trim();
    context.read<TestimonialsBloc>().add(GetTestimonialsRequested(
      search: searchQuery.isEmpty ? null : searchQuery,
      featured: _showFeaturedOnly,
      public: _showPublicOnly,
      rating: _selectedRating,
    ));
  }

  void _removeFilter(String filterType) {
    setState(() {
      if (filterType == 'featured') {
        _showFeaturedOnly = false;
      } else if (filterType == 'public') {
        _showPublicOnly = false;
      } else if (filterType == 'rating') {
        _selectedRating = null;
      }
    });
    _performSearch();
  }

  void _showFilterDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Filter Testimonials'),
        content: StatefulBuilder(
          builder: (context, setDialogState) {
            return Column(
              mainAxisSize: MainAxisSize.min,
              children: [
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
                const Text('Minimum Rating:', style: TextStyle(fontWeight: FontWeight.bold)),
                const SizedBox(height: 8),
                Row(
                  children: List.generate(5, (index) {
                    final rating = index + 1;
                    return Expanded(
                      child: RadioListTile<int>(
                        title: Text('$rating★+'),
                        value: rating,
                        groupValue: _selectedRating,
                        onChanged: (value) {
                          setDialogState(() {
                            _selectedRating = value;
                          });
                        },
                      ),
                    );
                  }),
                ),
                RadioListTile<int?>(
                  title: const Text('All Ratings'),
                  value: null,
                  groupValue: _selectedRating,
                  onChanged: (value) {
                    setDialogState(() {
                      _selectedRating = value;
                    });
                  },
                ),
              ],
            );
          },
        ),
        actions: [
          TextButton(
            onPressed: () {
              setState(() {
                _showFeaturedOnly = false;
                _showPublicOnly = false;
                _selectedRating = null;
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

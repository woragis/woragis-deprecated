import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../../domain/entities/ai_chat_entity.dart';
import '../bloc/money_bloc.dart';

class AiChatsListPage extends StatefulWidget {
  const AiChatsListPage({super.key});

  @override
  State<AiChatsListPage> createState() => _AiChatsListPageState();
}

class _AiChatsListPageState extends State<AiChatsListPage> {
  final _searchController = TextEditingController();
  bool _showArchivedOnly = false;
  String? _selectedAgent;

  @override
  void initState() {
    super.initState();
    // Load AI chats when page initializes
    _loadAiChats();
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _loadAiChats() {
    final searchQuery = _searchController.text.trim();
    context.read<MoneyBloc>().add(GetAiChatsRequested(
      search: searchQuery.isEmpty ? null : searchQuery,
      archived: _showArchivedOnly,
      agent: _selectedAgent,
    ));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('AI Chats'),
        backgroundColor: Colors.purple.shade600,
        foregroundColor: Colors.white,
        elevation: 0,
        actions: [
          IconButton(
            onPressed: () => _showFilterDialog(context),
            icon: const Icon(Icons.filter_list),
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
                hintText: 'Search chats...',
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
          if (_showArchivedOnly || _selectedAgent != null)
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              height: 50,
              child: ListView(
                scrollDirection: Axis.horizontal,
                children: [
                  if (_showArchivedOnly)
                    Chip(
                      label: const Text('Archived'),
                      onDeleted: () => _removeFilter('archived'),
                      deleteIcon: const Icon(Icons.close, size: 18),
                    ),
                  if (_selectedAgent != null)
                    Chip(
                      label: Text(_selectedAgent!),
                      onDeleted: () => _removeFilter('agent'),
                      deleteIcon: const Icon(Icons.close, size: 18),
                    ),
                ],
              ),
            ),

          // AI Chats list
          Expanded(
            child: BlocBuilder<MoneyBloc, MoneyState>(
              builder: (context, state) {
                if (state is MoneyLoading) {
                  return const Center(child: CircularProgressIndicator());
                } else if (state is MoneyError) {
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
                          'Error loading AI chats',
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
                            context.read<MoneyBloc>().add(GetAiChatsRequested());
                          },
                          child: const Text('Retry'),
                        ),
                      ],
                    ),
                  );
                } else if (state is AiChatsLoaded) {
                  if (state.aiChats.isEmpty) {
                    return Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(
                            Icons.chat_bubble_outline,
                            size: 64,
                            color: Colors.grey.shade400,
                          ),
                          const SizedBox(height: 16),
                          const Text(
                            'No AI chats found',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Text(
                            'Start a conversation with AI to explore your ideas!',
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
                      context.read<MoneyBloc>().add(GetAiChatsRequested());
                    },
                    child: ListView.builder(
                      padding: const EdgeInsets.all(16),
                      itemCount: state.aiChats.length,
                      itemBuilder: (context, index) {
                        final chat = state.aiChats[index];
                        return _buildChatCard(context, chat);
                      },
                    ),
                  );
                }

                return const SizedBox.shrink();
              },
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/money/ai-chats/create'),
        backgroundColor: Colors.purple.shade600,
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  Widget _buildChatCard(BuildContext context, AiChatEntity chat) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: () async {
          await context.push('/money/ai-chats/${chat.id}');
          // Reload list when returning from detail page
          _loadAiChats();
        },
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  CircleAvatar(
                    radius: 20,
                    backgroundColor: _getAgentColor(chat.agent),
                    child: Icon(
                      _getAgentIcon(chat.agent),
                      color: Colors.white,
                      size: 20,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          chat.title,
                          style: const TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        Text(
                          _getAgentDisplayName(chat.agent),
                          style: TextStyle(
                            fontSize: 12,
                            color: Colors.grey.shade600,
                          ),
                        ),
                      ],
                    ),
                  ),
                  if (chat.archived)
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 6,
                        vertical: 2,
                      ),
                      decoration: BoxDecoration(
                        color: Colors.grey.shade200,
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Text(
                        'Archived',
                        style: TextStyle(
                          fontSize: 10,
                          color: Colors.grey.shade600,
                        ),
                      ),
                    ),
                ],
              ),
              if (chat.messages.isNotEmpty) ...[
                const SizedBox(height: 12),
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: Colors.grey.shade50,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(
                    chat.messages.last.content,
                    style: TextStyle(
                      color: Colors.grey.shade700,
                      fontSize: 14,
                    ),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ],
              const SizedBox(height: 12),
              Row(
                children: [
                  Icon(
                    Icons.chat_bubble_outline,
                    size: 16,
                    color: Colors.grey.shade500,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    '${chat.messages.length} messages',
                    style: TextStyle(
                      fontSize: 12,
                      color: Colors.grey.shade500,
                    ),
                  ),
                  const Spacer(),
                  Text(
                    _formatDate(chat.createdAt),
                    style: TextStyle(
                      fontSize: 12,
                      color: Colors.grey.shade500,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Color _getAgentColor(AiAgent agent) {
    switch (agent) {
      case AiAgent.gpt4:
      case AiAgent.gpt4Turbo:
        return Colors.green;
      case AiAgent.gpt35Turbo:
        return Colors.blue;
      case AiAgent.claude3Opus:
      case AiAgent.claude3Sonnet:
      case AiAgent.claude3Haiku:
        return Colors.orange;
      case AiAgent.geminiPro:
        return Colors.purple;
      case AiAgent.custom:
        return Colors.grey;
    }
  }

  IconData _getAgentIcon(AiAgent agent) {
    switch (agent) {
      case AiAgent.gpt4:
      case AiAgent.gpt4Turbo:
      case AiAgent.gpt35Turbo:
        return Icons.smart_toy;
      case AiAgent.claude3Opus:
      case AiAgent.claude3Sonnet:
      case AiAgent.claude3Haiku:
        return Icons.psychology;
      case AiAgent.geminiPro:
        return Icons.auto_awesome;
      case AiAgent.custom:
        return Icons.settings;
    }
  }

  String _getAgentDisplayName(AiAgent agent) {
    switch (agent) {
      case AiAgent.gpt4:
        return 'GPT-4';
      case AiAgent.gpt4Turbo:
        return 'GPT-4 Turbo';
      case AiAgent.gpt35Turbo:
        return 'GPT-3.5 Turbo';
      case AiAgent.claude3Opus:
        return 'Claude 3 Opus';
      case AiAgent.claude3Sonnet:
        return 'Claude 3 Sonnet';
      case AiAgent.claude3Haiku:
        return 'Claude 3 Haiku';
      case AiAgent.geminiPro:
        return 'Gemini Pro';
      case AiAgent.custom:
        return 'Custom Model';
    }
  }

  void _performSearch() {
    final searchQuery = _searchController.text.trim();
    context.read<MoneyBloc>().add(GetAiChatsRequested(
      search: searchQuery.isEmpty ? null : searchQuery,
      archived: _showArchivedOnly,
      agent: _selectedAgent,
    ));
  }

  void _removeFilter(String filterType) {
    setState(() {
      if (filterType == 'archived') {
        _showArchivedOnly = false;
      } else if (filterType == 'agent') {
        _selectedAgent = null;
      }
    });
    _performSearch();
  }

  void _showFilterDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Filter AI Chats'),
        content: StatefulBuilder(
          builder: (context, setDialogState) {
            return Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                CheckboxListTile(
                  title: const Text('Archived only'),
                  value: _showArchivedOnly,
                  onChanged: (value) {
                    setDialogState(() {
                      _showArchivedOnly = value ?? false;
                    });
                  },
                ),
                const Divider(),
                const Text('Agent Filter:', style: TextStyle(fontWeight: FontWeight.bold)),
                ...AiAgent.values.map((agent) => RadioListTile<String>(
                  title: Text(_getAgentDisplayName(agent)),
                  value: agent.name,
                  groupValue: _selectedAgent,
                  onChanged: (value) {
                    setDialogState(() {
                      _selectedAgent = value;
                    });
                  },
                )),
                RadioListTile<String>(
                  title: const Text('All Agents'),
                  value: '',
                  groupValue: _selectedAgent,
                  onChanged: (value) {
                    setDialogState(() {
                      _selectedAgent = null;
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
                _showArchivedOnly = false;
                _selectedAgent = null;
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

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final difference = now.difference(date);

    if (difference.inDays == 0) {
      return 'Today';
    } else if (difference.inDays == 1) {
      return 'Yesterday';
    } else if (difference.inDays < 7) {
      return '${difference.inDays} days ago';
    } else {
      return '${date.day}/${date.month}/${date.year}';
    }
  }
}

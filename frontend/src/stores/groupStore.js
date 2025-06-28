import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useGroupStore = defineStore('group', () => {
  // State
  const groups = ref([])
  const userGroups = ref([])
  const publicGroups = ref([])
  const currentGroup = ref(null)
  const members = ref([])
  const posts = ref([])
  const events = ref([])
  const invitations = ref([])
  const requests = ref([])
  const loading = ref(false)
  const error = ref(null)

  // Getters
  const currentGroupMembers = computed(() => members.value)
  const currentGroupPosts = computed(() => posts.value)
  const currentGroupEvents = computed(() => events.value)
  const isCurrentUserMember = computed(() => {
    return currentGroup.value?.member || false
  })
  const isCurrentUserAdmin = computed(() => {
    return currentGroup.value?.administrator || false
  })

  // Actions
  const fetchGroups = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await fetch('http://localhost:8081/allGroups', {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        groups.value = data.groups || []
      } else {
        throw new Error('Failed to fetch groups')
      }
    } catch (err) {
      error.value = err.message
      console.error('Error fetching groups:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchUserGroups = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await fetch('http://localhost:8081/userGroups', {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        groups.value = data.groups || []
      } else {
        throw new Error('Failed to fetch user groups')
      }
    } catch (err) {
      error.value = err.message
      console.error('Error fetching user groups:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchGroupDetails = async (groupId) => {
    loading.value = true
    error.value = null
    try {
      const response = await fetch(`http://localhost:8081/groupInfo?groupId=${groupId}`, {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        currentGroup.value = data.groups && data.groups.length > 0 ? data.groups[0] : null
      } else {
        throw new Error('Failed to fetch group details')
      }
    } catch (err) {
      error.value = err.message
      console.error('Error fetching group details:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchGroupMembers = async (groupId) => {
    try {
      const response = await fetch(`http://localhost:8081/groupMembers?groupId=${groupId}`, {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        members.value = data.users || []
      }
    } catch (err) {
      console.error('Error fetching group members:', err)
    }
  }

  const fetchGroupPosts = async (groupId) => {
    try {
      const response = await fetch(`http://localhost:8081/groupPosts?groupId=${groupId}`, {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        posts.value = data.posts || []
      }
    } catch (err) {
      console.error('Error fetching group posts:', err)
    }
  }

  const fetchGroupEvents = async (groupId) => {
    try {
      const response = await fetch(`http://localhost:8081/groupEvents?groupId=${groupId}`, {
        credentials: 'include'
      })
      if (response.ok) {
        const data = await response.json()
        events.value = data.events || []
      }
    } catch (err) {
      console.error('Error fetching group events:', err)
    }
  }

  const joinGroup = async (groupId, isPublic = true) => {
    const endpoint = isPublic 
      ? 'http://localhost:8081/joinPublicGroup'
      : 'http://localhost:8081/newGroupRequest'
    
    try {
      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ groupId })
      })

      if (response.ok && isPublic) {
        // Update current group membership status
        if (currentGroup.value && currentGroup.value.id === groupId) {
          currentGroup.value.member = true
        }
        // Refresh group data
        await fetchGroupDetails(groupId)
        await fetchGroupMembers(groupId)
        return true
      }
      return response.ok
    } catch (err) {
      console.error('Error joining group:', err)
      return false
    }
  }

  const leaveGroup = async (groupId) => {
    try {
      const response = await fetch('http://localhost:8081/leaveGroup', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({ groupId })
      })

      if (response.ok) {
        // Update current group membership status
        if (currentGroup.value && currentGroup.value.id === groupId) {
          currentGroup.value.member = false
        }
        // Refresh group data
        await fetchGroupDetails(groupId)
        await fetchGroupMembers(groupId)
        return true
      }
      return false
    } catch (err) {
      console.error('Error leaving group:', err)
      return false
    }
  }

  const clearCurrentGroup = () => {
    currentGroup.value = null
    members.value = []
    posts.value = []
    events.value = []
  }

  return {
    // State
    groups,
    currentGroup,
    members,
    posts,
    events,
    loading,
    error,
    
    // Getters
    currentGroupMembers,
    currentGroupPosts,
    currentGroupEvents,
    isCurrentUserMember,
    isCurrentUserAdmin,
    
    // Actions
    fetchGroups,
    fetchUserGroups,
    fetchGroupDetails,
    fetchGroupMembers,
    fetchGroupPosts,
    fetchGroupEvents,
    joinGroup,
    leaveGroup,
    clearCurrentGroup
  }
})

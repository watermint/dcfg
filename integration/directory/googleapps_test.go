package directory

import (
	"fmt"
	"google.golang.org/api/admin/directory/v1"
	"testing"
)

func TestGoogleAppsWithCache_Users(t *testing.T) {
	validate := func(expected, results []*admin.User) {
		if len(results) != len(expected) {
			t.Errorf("Invalid result: %v", results)
		}
		for _, x := range expected {
			found := false
			for _, y := range results {
				if x.PrimaryEmail == y.PrimaryEmail {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("User not found: %v", x)
			}
		}
	}

	for numUsers := 0; numUsers < 5; numUsers++ {
		expectedUsers := make([]*admin.User, 0, numUsers)
		for j := 0; j < numUsers; j++ {
			expectedUsers = append(expectedUsers, &admin.User{
				PrimaryEmail: fmt.Sprintf("a%d@example.com", j),
			})
		}

		mock := GoogleAppsMock{
			MockUsers: expectedUsers,
		}
		cache := GoogleAppsWithCache{
			Resolver: &mock,
		}

		mock.Preload()
		cache.Preload()

		validate(expectedUsers, cache.Users())
	}
}

func TestGoogleAppsMock_Groups(t *testing.T) {
	validate := func(expected, results []*admin.Group) {
		if len(results) != len(expected) {
			t.Errorf("Invalid result: %v", results)
		}
		for _, x := range expected {
			found := false
			for _, y := range results {
				if x.Email == y.Email {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Group not found: %v", x)
			}
		}
	}

	for numGroups := 0; numGroups < 5; numGroups++ {
		expectedGroups := make([]*admin.Group, 0, numGroups)
		for j := 0; j < numGroups; j++ {
			expectedGroups = append(expectedGroups, &admin.Group{
				Email: fmt.Sprintf("a%d@example.com", j),
			})
		}

		mock := GoogleAppsMock{
			MockGroups: expectedGroups,
		}
		cache := GoogleAppsWithCache{
			Resolver: &mock,
		}

		mock.Preload()
		cache.Preload()

		validate(expectedGroups, cache.Groups())
	}
}

func TestGoogleAppsMock_CustomerUsers(t *testing.T) {
	validate := func(expected, results []*admin.User) {
		if len(results) != len(expected) {
			t.Errorf("Invalid result: %v", results)
		}
		for _, x := range expected {
			found := false
			for _, y := range results {
				if x.PrimaryEmail == y.PrimaryEmail {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("User not found: %v", x)
			}
		}
	}

	for numUsers := 0; numUsers < 5; numUsers++ {
		expectedUsers := make([]*admin.User, 0, numUsers)
		for j := 0; j < numUsers; j++ {
			expectedUsers = append(expectedUsers, &admin.User{
				PrimaryEmail: fmt.Sprintf("a%d@example.com", j),
			})
		}
		mock := GoogleAppsMock{
			MockCustomers: map[string][]*admin.User{
				"mock_customer": expectedUsers,
				"empty":         []*admin.User{},
			},
		}
		cache := GoogleAppsWithCache{
			Resolver: &mock,
		}

		mock.Preload()
		cache.Preload()

		validate([]*admin.User{}, cache.CustomerUsers("empty"))
		validate([]*admin.User{}, cache.CustomerUsers("no_existent"))
		validate(expectedUsers, cache.CustomerUsers("mock_customer"))
	}
}

func TestGoogleAppsWithCache_GroupMembers(t *testing.T) {
	validate := func(expected, results []*admin.Member) {
		if len(expected) != len(results) {
			t.Errorf("Invalid result: %v", results)
		}
		for _, x := range expected {
			found := false
			for _, y := range results {
				if x.Email == y.Email && x.Type == y.Type {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Member not found: %v", x)
			}
		}
	}

	for numMembers := 0; numMembers < 5; numMembers++ {
		expectedMembers := make([]*admin.Member, 0, numMembers)
		for j := 0; j < numMembers; j++ {
			expectedMembers = append(expectedMembers, &admin.Member{
				Type:  "USER",
				Email: fmt.Sprintf("a%d@example.com", j),
			})
		}
		mock := GoogleAppsMock{
			MockMembers: map[string][]*admin.Member{
				"mock_group":  expectedMembers,
				"empty_group": []*admin.Member{},
			},
		}
		cache := GoogleAppsWithCache{
			Resolver: &mock,
		}

		mock.Preload()
		cache.Preload()

		validate([]*admin.Member{}, cache.GroupMembers("no_existent"))
		validate([]*admin.Member{}, cache.GroupMembers("empty_group"))
		validate(expectedMembers, cache.GroupMembers("mock_group"))
	}
}

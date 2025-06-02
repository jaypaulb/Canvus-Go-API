package canvus

import (
	"context"
	"testing"
	"time"
)

func TestFolderLifecycle(t *testing.T) {
	ctx := context.Background()
	admin, _, err := getTestAdminClientFromSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}

	// Create a root test folder for isolation
	rootFolderName := "testfolder_root_" + time.Now().Format("20060102150405")
	rootFolder, err := admin.CreateFolder(ctx, CreateFolderRequest{Name: rootFolderName})
	if err != nil {
		t.Fatalf("failed to create root test folder: %v", err)
	}
	// Clean up: delete all children, then the root folder itself
	defer func() {
		_ = admin.DeleteFolderContents(ctx, rootFolder.ID)
		_ = admin.DeleteFolder(ctx, rootFolder.ID)
	}()

	// Create a folder inside the root test folder
	folderName := rootFolderName + "_child"
	folder, err := admin.CreateFolder(ctx, CreateFolderRequest{Name: folderName, ParentID: rootFolder.ID})
	if err != nil {
		t.Fatalf("failed to create folder: %v", err)
	}
	// Clean up: delete this folder at the end
	defer func() { _ = admin.DeleteFolder(ctx, folder.ID) }()

	// List folders and check the new folder is present
	folders, err := admin.ListFolders(ctx)
	if err != nil {
		t.Errorf("failed to list folders: %v", err)
	}
	found := false
	for _, f := range folders {
		if f.ID == folder.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created folder not found in list")
	}

	// Get the folder by ID
	got, err := admin.GetFolder(ctx, folder.ID)
	if err != nil {
		t.Errorf("failed to get folder: %v", err)
	}
	if got.Name != folderName {
		t.Errorf("expected folder name %q, got %q", folderName, got.Name)
	}

	// Rename the folder
	newName := folderName + "_renamed"
	renamed, err := admin.RenameFolder(ctx, folder.ID, newName)
	if err != nil {
		t.Errorf("failed to rename folder: %v", err)
	}
	if renamed.Name != newName {
		t.Errorf("expected renamed folder name %q, got %q", newName, renamed.Name)
	}

	// Create a subfolder inside the root test folder
	subName := rootFolderName + "_sub"
	sub, err := admin.CreateFolder(ctx, CreateFolderRequest{Name: subName, ParentID: rootFolder.ID})
	if err != nil {
		t.Fatalf("failed to create subfolder: %v", err)
	}
	defer func() { _ = admin.DeleteFolder(ctx, sub.ID) }()

	// Move subfolder to the child folder
	moved, err := admin.MoveFolder(ctx, sub.ID, folder.ID, "replace")
	if err != nil {
		t.Errorf("failed to move subfolder: %v", err)
	} else if moved.ParentID != folder.ID {
		t.Errorf("expected subfolder parent to be %q, got %q", folder.ID, moved.ParentID)
	}

	// Copy folder (copy subfolder into root test folder)
	copied, err := admin.CopyFolder(ctx, sub.ID, rootFolder.ID, "replace")
	if err != nil {
		t.Errorf("failed to copy folder: %v", err)
	}
	defer func() { _ = admin.DeleteFolder(ctx, copied.ID) }()

	// Trash the copied folder
	trashed, err := admin.TrashFolder(ctx, copied.ID, "trash.1000")
	if err != nil {
		t.Errorf("failed to trash folder: %v", err)
	}
	if !trashed.InTrash {
		t.Errorf("expected folder to be in trash, got in_trash=%v", trashed.InTrash)
	}

	// Delete all children of the root test folder (should not error)
	err = admin.DeleteFolderContents(ctx, rootFolder.ID)
	if err != nil {
		t.Errorf("failed to delete folder contents: %v", err)
	}

	// Permissions: get and set on the root test folder
	perms, err := admin.GetFolderPermissions(ctx, rootFolder.ID)
	if err != nil {
		t.Errorf("failed to get folder permissions: %v", err)
	}
	perms.EditorsCanShare = false
	updated, err := admin.SetFolderPermissions(ctx, rootFolder.ID, *perms)
	if err != nil {
		t.Errorf("failed to set folder permissions: %v", err)
	}
	if updated.EditorsCanShare != false {
		t.Errorf("expected EditorsCanShare to be false, got %v", updated.EditorsCanShare)
	}
}

func TestFolderInvalidCases(t *testing.T) {
	ctx := context.Background()
	admin, _, err := getTestAdminClientFromSettings()
	if err != nil {
		t.Fatalf("failed to load test settings: %v", err)
	}

	// Get non-existent folder
	_, err = admin.GetFolder(ctx, "nonexistent-folder-id")
	if err == nil {
		t.Errorf("expected error for non-existent folder, got nil")
	}

	// Delete non-existent folder
	err = admin.DeleteFolder(ctx, "nonexistent-folder-id")
	if err == nil {
		t.Errorf("expected error for deleting non-existent folder, got nil")
	}

	// Rename non-existent folder
	_, err = admin.RenameFolder(ctx, "nonexistent-folder-id", "newname")
	if err == nil {
		t.Errorf("expected error for renaming non-existent folder, got nil")
	}

	// Move non-existent folder
	_, err = admin.MoveFolder(ctx, "nonexistent-folder-id", "", "replace")
	if err == nil {
		t.Errorf("expected error for moving non-existent folder, got nil")
	}

	// Copy non-existent folder
	_, err = admin.CopyFolder(ctx, "nonexistent-folder-id", "", "replace")
	if err == nil {
		t.Errorf("expected error for copying non-existent folder, got nil")
	}

	// Trash non-existent folder
	_, err = admin.TrashFolder(ctx, "nonexistent-folder-id", "trash.1000")
	if err == nil {
		t.Errorf("expected error for trashing non-existent folder, got nil")
	}

	// Delete contents of non-existent folder
	err = admin.DeleteFolderContents(ctx, "nonexistent-folder-id")
	if err == nil {
		t.Errorf("expected error for deleting contents of non-existent folder, got nil")
	}

	// Get permissions of non-existent folder
	_, err = admin.GetFolderPermissions(ctx, "nonexistent-folder-id")
	if err == nil {
		t.Errorf("expected error for getting permissions of non-existent folder, got nil")
	}

	// Set permissions of non-existent folder
	perms := FolderPermissions{EditorsCanShare: false}
	_, err = admin.SetFolderPermissions(ctx, "nonexistent-folder-id", perms)
	if err == nil {
		t.Errorf("expected error for setting permissions of non-existent folder, got nil")
	}
}

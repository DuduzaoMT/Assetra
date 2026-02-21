import { useState, useEffect, useCallback } from "react";
import { userManagementService } from "../services/userService";
import { toast } from "react-hot-toast";

export function useUserManagement() {
  const [usersList, setUsersList] = useState([]);
  const [isUsersLoading, setIsUsersLoading] = useState(false);
  const [usersError, setUsersError] = useState(null);

  const fetchAllUsers = async () => {
    try {
      setIsUsersLoading(true);
      setUsersError(null);
      const usersData = await userManagementService.getAllUsers();
      setUsersList(usersData);
    } catch (fetchError) {
      setUsersError(fetchError.message);
      toast.error(fetchError.message);
    } finally {
      setIsUsersLoading(false);
    }
  };

  const updateUserInformation = async (userId, userUpdateData) => {
    try {
      const updatedUserData = await userManagementService.updateUserData(
        userId,
        userUpdateData
      );
      setUsersList((currentUsers) =>
        currentUsers.map((user) =>
          user.id === userId ? updatedUserData : user
        )
      );
      toast.success("User updated successfully!");
      return { success: true };
    } catch (updateError) {
      toast.error(updateError.message);
      return { success: false, error: updateError.message };
    }
  };

  const removeUserFromSystem = async (userId) => {
    try {
      await userManagementService.deleteUserById(userId);
      setUsersList((currentUsers) =>
        currentUsers.filter((user) => user.id !== userId)
      );
      toast.success("User deleted successfully!");
      return { success: true };
    } catch (deleteError) {
      toast.error(deleteError.message);
      return { success: false, error: deleteError.message };
    }
  };

  useEffect(() => {
    fetchAllUsers();
  }, []);

  return {
    usersList,
    isUsersLoading,
    usersError,
    refreshUsersList: fetchAllUsers,
    updateUserInformation,
    removeUserFromSystem,
  };
}

export function useIndividualUser(userId) {
  const [userData, setUserData] = useState(null);
  const [isUserLoading, setIsUserLoading] = useState(false);
  const [userError, setUserError] = useState(null);

  const fetchUserData = useCallback(async () => {
    if (!userId) return;

    try {
      setIsUserLoading(true);
      setUserError(null);
      const fetchedUserData = await userManagementService.getUserById(userId);
      setUserData(fetchedUserData);
    } catch (fetchError) {
      setUserError(fetchError.message);
      toast.error(fetchError.message);
    } finally {
      setIsUserLoading(false);
    }
  }, [userId]);

  useEffect(() => {
    fetchUserData();
  }, [fetchUserData]);

  return {
    userData,
    isUserLoading,
    userError,
    refreshUserData: fetchUserData,
  };
}

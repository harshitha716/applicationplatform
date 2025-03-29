import React from 'react';
import { useSelector } from 'react-redux';
import { useGetOrganizationMembershipRequestsAllQuery } from 'apis/people';
import { ZAMP_ICON } from 'constants/icons';
import { useLogout } from 'hooks/useLogout';
import Image from 'next/image';
import { RootState } from 'store';
import { MembershipRequested } from 'components/MembershipRequested';

const OrgMembershipPending = () => {
  const userEmail = useSelector((state: RootState) => state?.user?.user)?.user_email;
  const { logout } = useLogout();
  const { data: membershipRequests, isLoading: isLoadingMembershipRequests } =
    useGetOrganizationMembershipRequestsAllQuery();

  const logoutButton = {
    text: 'Logout',
    onClick: logout,
  };

  // TODO: Loading animation
  if (isLoadingMembershipRequests) {
    return (
      <div className='w-screen h-screen flex flex-col bg-white justify-center items-center'>
        <Image
          width={60}
          height={60}
          alt='zamp logo'
          className='w-8 align-middle cursor-pointer'
          src={ZAMP_ICON}
          priority={true}
        />
      </div>
    );
  }

  if (membershipRequests && membershipRequests?.length > 0) {
    return (
      <MembershipRequested
        text='Your account is pending approval'
        subText='We have notified the organization admin. You will receive an email when your membership request is approved.'
        userEmail={userEmail || ''}
        actionItems={[logoutButton]}
      />
    );
  }

  return (
    <MembershipRequested
      text='Thank you for your interest in Zamp'
      subText='We have received your signup request and our team will review it shortly.'
      userEmail={userEmail || ''}
      actionItems={[logoutButton]}
    />
  );
};

export default OrgMembershipPending;

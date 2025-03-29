import { useSelector } from 'react-redux';
import { render, screen } from '@testing-library/react';
import { useGetOrganizationMembershipRequestsAllQuery } from 'apis/people';
import { useLogout } from 'hooks/useLogout';
import OrgMembershipPending from 'modules/login/OrgMembershipPending';

jest.mock('next/font/google', () => ({
  Inter: jest.fn(() => ({
    className: 'mocked-inter',
    variable: '--font-inter',
  })),
}));

jest.mock('react-redux', () => ({
  useSelector: jest.fn(),
}));

jest.mock('apis/people', () => ({
  useGetOrganizationMembershipRequestsAllQuery: jest.fn(),
}));

jest.mock('hooks/useLogout', () => ({
  useLogout: jest.fn(),
}));

jest.mock('next/image', () => ({
  __esModule: true,
  default: (props: any) => <img {...props} />,
}));

describe('OrgMembershipPending', () => {
  const mockLogout = jest.fn();
  const mockUserEmail = 'test@example.com';

  beforeEach(() => {
    (useSelector as unknown as jest.Mock).mockReturnValue({ user_email: mockUserEmail });
    (useLogout as jest.Mock).mockReturnValue({ logout: mockLogout });
  });

  it('should show loading state', () => {
    (useGetOrganizationMembershipRequestsAllQuery as jest.Mock).mockReturnValue({
      isLoading: true,
      data: null,
    });

    render(<OrgMembershipPending />);

    expect(screen.getByAltText('zamp logo')).toBeInTheDocument();
  });

  it('should show pending approval message when membership requests exist', () => {
    (useGetOrganizationMembershipRequestsAllQuery as jest.Mock).mockReturnValue({
      isLoading: false,
      data: [{ id: 1 }],
    });

    render(<OrgMembershipPending />);

    expect(screen.getByText('Your account is pending approval')).toBeInTheDocument();
    expect(
      screen.getByText(
        'We have notified the organization admin. You will receive an email when your membership request is approved.',
      ),
    ).toBeInTheDocument();
    expect(screen.getByText(mockUserEmail)).toBeInTheDocument();
    expect(screen.getByText('Logout')).toBeInTheDocument();
  });

  it('should show signup request message when no membership requests exist', () => {
    (useGetOrganizationMembershipRequestsAllQuery as jest.Mock).mockReturnValue({
      isLoading: false,
      data: [],
    });

    render(<OrgMembershipPending />);

    expect(screen.getByText('Thank you for your interest in Zamp')).toBeInTheDocument();
    expect(
      screen.getByText('We have received your signup request and our team will review it shortly.'),
    ).toBeInTheDocument();
    expect(screen.getByText(mockUserEmail)).toBeInTheDocument();
    expect(screen.getByText('Logout')).toBeInTheDocument();
  });
});

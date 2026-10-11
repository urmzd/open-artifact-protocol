<oap:target id="order-model-file">
class Order < ApplicationRecord
  <oap:target id="order-associations">
    belongs_to <oap:target id="user-association">:user</oap:target>
    has_many <oap:target id="order-items-association">:order_items</oap:target>
    has_one <oap:target id="shipping-address-association">:shipping_address</oap:target>
    has_many <oap:target id="payments-association">:payments</oap:target>
  </oap:target>

  <oap:target id="order-enum-and-state">
    enum <oap:target id="status-attribute">status</oap:target>: {
      <oap:target id="status-pending">pending</oap:target>: 0,
      <oap:target id="status-confirmed">confirmed</oap:target>: 1,
      <oap:target id="status-shipped">shipped</oap:target>: 2,
      <oap:target id="status-delivered">delivered</oap:target>: 3,
      <oap:target id="status-cancelled">cancelled</oap:target>: 4
    }
  </oap:target>

  <oap:target id="order-validations">
    validates <oap:target id="order-status-presence">:status</oap:target>, presence: true
    validates <oap:target id="order-total-presence">:total_cents</oap:target>, presence: true, numericality: { greater_than_or_equal_to: <oap:target id="minimum-order-total">0</oap:target> }
    validates <oap:target id="order-subtotal-presence">:subtotal_cents</oap:target>, presence: true, numericality: { greater_than_or_equal_to: <oap:target id="minimum-order-subtotal">0</oap:target> }
    validates <oap:target id="order-tax-presence">:tax_cents</oap:target>, presence: true, numericality: { greater_than_or_equal_to: <oap:target id="minimum-order-tax">0</oap:target> }

    validate <oap:target id="validate-status-transitions-method">:validate_status_transitions</oap:target>
    validate <oap:target id="validate-shipment-eligibility-method">:validate_shipment_eligibility</oap:target>
  </oap:target>

  <oap:target id="order-scopes">
    scope <oap:target id="recent-scope-name">:recent</oap:target>, -> { order(<oap:target id="recent-order-column">created_at</oap:target>: <oap:target id="recent-order-direction">:desc</oap:target>) }
    scope <oap:target id="by-status-scope-name">:by_status</oap:target>, ->(status) { where(<oap:target id="by-status-column">status</oap:target>: status) }
    scope <oap:target id="by-date-range-scope-name">:by_date_range</oap:target>, ->(start_date, end_date) { where(<oap:target id="date-range-column">created_at</oap:target>: start_date..end_date) }
    scope <oap:target id="high-value-scope-name">:high_value</oap:target>, -> { where(<oap:target id="high-value-column">total_cents</oap:target> > <oap:target id="high-value-threshold">10_000</oap:target>) }
    scope <oap:target id="pending-shipment-scope-name">:pending_shipment</oap:target>, -> { where(<oap:target id="pending-shipment-status-column">status</oap:target>: <oap:target id="pending-shipment-status-value">:confirmed</oap:target>) }
  </oap:target>

  <oap:target id="order-callbacks">
    before_create <oap:target id="before-create-callback-method">:set_initial_status</oap:target>
    after_update <oap:target id="after-update-callback-method">:track_status_change</oap:target>
  </oap:target>

  <oap:target id="order-instance-methods">
    def <oap:target id="subtotal-method-name">subtotal</oap:target>
      <oap:target id="subtotal-method-body">subtotal_cents.to_d / 100</oap:target>
    end

    def <oap:target id="tax-amount-method-name">tax_amount</oap:target>
      <oap:target id="tax-amount-method-body">tax_cents.to_d / 100</oap:target>
    end

    def <oap:target id="total-method-name">total</oap:target>
      <oap:target id="total-method-body">total_cents.to_d / 100</oap:target>
    end

    def <oap:target id="apply-coupon-method-name">apply_coupon</oap:target>(coupon_code)
      <oap:target id="apply-coupon-method-body"># Apply coupon logic here</oap:target>
    end

    def <oap:target id="can-cancel-method-name">can_cancel?</oap:target>
      <oap:target id="can-cancel-method-body">pending? || confirmed?</oap:target>
    end

    def <oap:target id="ship-method-name">ship!</oap:target>
      <oap:target id="ship-method-body">update!(status: :shipped)</oap:target>
    end

    def <oap:target id="complete-method-name">complete!</oap:target>
      <oap:target id="complete-method-body">update!(status: :delivered)</oap:target>
    end
  </oap:target>

  <oap:target id="order-status-machine">
    <oap:target id="status-transition-pending-confirmed">AASM-like: pending -> confirmed</oap:target>
    <oap:target id="status-transition-confirmed-shipped">AASM-like: confirmed -> shipped</oap:target>
    <oap:target id="status-transition-shipped-delivered">AASM-like: shipped -> delivered</oap:target>
    <oap:target id="status-transition-cancel">AASM-like: any allowed -> cancelled</oap:target>
  </oap:target>

  private

  <oap:target id="order-private-methods">
    def <oap:target id="set-initial-status-method-name">set_initial_status</oap:target>
      <oap:target id="set-initial-status-method-body">self.status ||= :pending</oap:target>
    end

    def <oap:target id="track-status-change-method-name">track_status_change</oap:target>
      <oap:target id="track-status-change-method-body"># Track order status changes here</oap:target>
    end

    def <oap:target id="validate-status-transitions-method-name">validate_status_transitions</oap:target>
      <oap:target id="validate-status-transitions-method-body"># Validate allowed status transitions here</oap:target>
    end

    def <oap:target id="validate-shipment-eligibility-method-name">validate_shipment_eligibility</oap:target>
      <oap:target id="validate-shipment-eligibility-method-body"># Ensure the order can be shipped before transitioning</oap:target>
    end
  end
end
</oap:target>
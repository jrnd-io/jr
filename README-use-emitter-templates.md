# Emitter usage guide

This is a list of `jr` commands to run the templates that are linked togheter using emitters.

The preconfigured emitters examples have Kafka as output.

# Global templates 

These templates:

* util_userid.tpl
* util_ip.tpl

Are used to preload lists of values used in the other templates in the emitter configuration.

## Shoes emitter

    jr emitter run shoe shoe_clickstream shoe_customer shoe_order 

## Clickstream emitter

    jr emitter run clickstream_codes_schema clickstream_schema clickstream_users_schema

## Credit cards and Transactions emitter

    jr emitter run util_userid credit_cards transactions

## Gaming emitter

    jr emitter run gaming_games gaming_players gaming_player_activity

## Insurrance emitter

    jr emitter run insurance_customers insurance_customer_activity  insurance_offers

## Payroll emitter

    jr emitter run payroll_employee payroll_bonus payroll_employee_location

## Pizza Orders emitter

    jr emitter run  pizza_store_util pizza_orders pizza_orders_completed pizza_orders_cancelled

## Pageviews emitter

    jr emitter run util_userid pageviews_schema

## Rating emitter

    jr emitter run util_userid rating_schema    

## Stock trade emitter

    jr emitter run util_userid stock_trades-schema    

## Siem Logs emitter

    jr emitter run util_ip siem_logs

## Inventory Product emitter

    jr emitter run product inventory